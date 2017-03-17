package stomp

import (
	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
)

// NodeStomp is a class with methods for STOMP comands.
// It inherits from transport.Node.
type NodeStomp struct {
	*transport.Node
}

// NewNode creates a new NodeStomp object and returns it.
func NewNode() *NodeStomp {
	n := &NodeStomp{transport.NewNode()}

	stompProcessor := &ProcessorStomp{node: n}
	n.AddCmdProcessor(stompProcessor)

	return n
}

// SendFrame sends a frame to ActiveLink with certain ID.
func (n *NodeStomp) SendFrame(topic string, frame *frame.Frame) {

	log.WithField("topic=", topic).Debugf("funcSendFrame() enter")

	n.CommandCh <- &CommandSendFrameStomp{
		transport.NodeCommand{Cmd: stompSendFrameCommand},
		*frame,
		topic,
	}
}

/*
// RecieveFrame
func (n *NodeStomp) RecieveFrame(activeLinkID string, frame *frame.Frame) {

	log.Debugf("funcSendFrame() for [%s]", activeLinkID)

	n.CommandCh <- &CommandSendFrameStomp{
		transport.NodeCommand{Cmd: stompSendFrameCommand},
		*frame,
		activeLinkID,
	}
}
*/
