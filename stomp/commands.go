package stomp

import (
	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
)

type StompCommandSendFrame struct {
	transport.NodeCommand
	frame        frame.Frame
	linkActiveID string
}

const (
	stompSendFrameCommand transport.CommandID = 100
)
