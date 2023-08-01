package packet

import (
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio/frame"
)

type Frame struct {
	FType frame.Type
	Data  []byte
}

type Packet struct {
	FType frame.Type
	PType Type
	Data  []byte
}
