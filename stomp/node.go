package stomp

import (
	"errors"

	"github.com/KristinaEtc/bdmq/transport"
)

type subscripted struct {
	ok  bool
	err error
}

// NodeStomp is a class with methods for STOMP comands.
// It inherits from transport.Node.
type NodeStomp struct {
	*transport.Node
	subscriptions map[string]map[string]*chan transport.Frame // groupped by topic, then by LinkActiveID
}

// NewNode creates a new NodeStomp object and returns it.
func NewNode() *NodeStomp {

	n := &NodeStomp{transport.NewNode(),
		make(map[string]map[string]*chan transport.Frame)}

	return n
}

// SendFrame sends a frame to ActiveLink with certain ID.
func (n *NodeStomp) SendFrame(linkActiveID string, frame *transport.Frame) {

	log.WithField("linkActiveID=", linkActiveID).Debugf("funcSendFrame() enter")

	n.CommandCh <- &CommandSendFrameStomp{
		transport.NodeCommand{Cmd: stompSendFrameCommand},
		*frame,
		linkActiveID,
	}
}

//Subscribe sends a frame to subscribe activeLink with ID = ActiveLinkID with topic.
func (n *NodeStomp) Subscribe(topic string) (*chan transport.Frame, error) {

	_, ok := n.LinkActives[linkActiveID]
	if !ok {
		log.Errorf("Error: no link actives with such ID=[%s]", linkActiveID)
		log.Debugf("have %+v", n.LinkActives)
		return nil, errors.New("Error: no link actives with such ID")
	}

	var sub = make(chan subscripted, 0)

	n.CommandCh <- &CommandSubscribeStomp{
		transport.NodeCommand{Cmd: stompSubscribeCommand},
		topic,
		linkActiveID,
		sub,
	}

	select {
	case subscriptionCompleted := <-sub:
		{
			if subscriptionCompleted.ok {
				return n.subscriptions[topic][linkActiveID], nil
			}
			log.Errorf("Could not subscribe: aclive Link=[%s], topic=[%s]: %s", linkActiveID, topic, subscriptionCompleted.err.Error())
			return nil, subscriptionCompleted.err
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
