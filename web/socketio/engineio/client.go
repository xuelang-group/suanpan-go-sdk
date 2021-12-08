package engineio

import (
	"encoding/json"
	"fmt"
	"io"
)

func ReadHandshake(r io.Reader) (*Session, error) {
	p, err := NewDecoder(r).Decode()
	if err != nil {
		return nil, fmt.Errorf("decode initial engine.io packet: %w", err)
	}
	if p.Type != OPEN {
		return nil, fmt.Errorf("unexpected engine.io packet type(expected=%v, got=%v)", OPEN, p.Type)
	}

	var session Session
	if err := json.NewDecoder(p.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("invalid session json: %w", err)
	}
	return &session, nil
}

func MessagePrefix() []byte {
	return Prefix(MESSAGE)
}

func Prefix(packetType PacketType) []byte {
	return []byte{byte(packetType) + '0'}
}
