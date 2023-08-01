package session

import (
	"github.com/xuelang-group/suanpan-go-sdk/web/socketio/engineio/frame"
)

// FrameType is type of message frame.
type FrameType frame.Type

const (
	// TEXT is text type message.
	TEXT = FrameType(frame.String)
	// BINARY is binary type message.
	BINARY = FrameType(frame.Binary)
)
