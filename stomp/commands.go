package stomp

import (
	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
)

// CommandSendFrameStomp is a struct with will send to StompNode proseccor
type CommandSendFrameStomp struct {
	transport.NodeCommand
	frame        frame.Frame
	linkActiveID string
}

const (
	stompSendFrameCommand transport.CommandID = 100
)
