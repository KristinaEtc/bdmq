package stomp

import (
	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
)

type StompNode struct {
	*transport.Node
}

func NewNode() *StompNode {
	n := &StompNode{transport.NewNode()}

	StompProcesser := &StompProcesser{node: n}
	n.AddCmdProcessor(StompProcesser)

	return n
}

func (n *StompNode) SendFrame(activeLinkID string, frame *frame.Frame) {

	log.Warnf("funcSendFrame() for [%s]", activeLinkID)

	n.CommandCh <- &StompCommandSendFrame{
		transport.NodeCommand{Cmd: stompSendFrameCommand},
		*frame,
		activeLinkID,
	}
}

func (n *StompNode) RecieveFrame(activeLinkID string, frame *frame.Frame) {

	log.Debugf("funcSendFrame() for [%s]", activeLinkID)

	n.CommandCh <- &StompCommandSendFrame{
		transport.NodeCommand{Cmd: stompSendFrameCommand},
		*frame,
		activeLinkID,
	}
}
