package stomp

import (
	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
)

// Nodes for using them in handlers!
var nodes = make(map[string]*NodeStomp)

type subscripted struct {
	ok bool
}

// Subscription is a struct with parametrs for sending by topics
type Subscription struct {
	//	linkActiveID string
	ch chan frame.Frame
}

// NodeStomp is a class with methods for STOMP comands.
// It inherits from transport.Node.
type NodeStomp struct {
	*transport.Node
	subscriptions map[string]Subscription // groupped by topic
}

// NewNode creates a new NodeStomp object and returns it.
func NewNode() *NodeStomp {

	n := &NodeStomp{transport.NewNode(),
		make(map[string]Subscription)}

	// TODO; make it ok; now works with one node only
	n.NodeID = "NodeStomp-ID"

	nodes[n.NodeID] = n

	return n
}

// SendFrame sends a frame to ActiveLink with certain ID.
func (n *NodeStomp) SendFrame(topic string, frame frame.Frame) {

	log.WithField("topic", topic).Debugf("funcSendFrame()")

	n.CommandCh <- &CommandSendFrameStomp{
		transport.NodeCommand{Cmd: stompSendFrameCommand},
		frame,
		topic,
	}
}

// ReceiveFrame used in handlers for recieving a frame.
func (n *NodeStomp) ReceiveFrame(linkActiveID string, topic string, frame frame.Frame) {

	log.WithField("linkActiveID ", linkActiveID).Debugf("func ReceiveFrame()")

	n.CommandCh <- &CommandReceiveFrameStomp{
		transport.NodeCommand{Cmd: stompReceiveFrameCommand},
		frame,
		linkActiveID,
		topic,
	}
}

//Subscribe sends a frame to subscribe activeLink with ID = ActiveLinkID with topic.
func (n *NodeStomp) Subscribe(topic string) (chan frame.Frame, error) {

	var sub = make(chan subscripted, 0)

	n.CommandCh <- &CommandSubscribeStomp{
		transport.NodeCommand{Cmd: stompSubscribeCommand},
		topic,
		sub,
	}

	select {
	case _ = <-sub:
		{
			return n.subscriptions[topic].ch, nil
		}
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
