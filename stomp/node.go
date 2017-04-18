package stomp

import (
	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
)

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
	handlers      []*HandlerStomp
}

// NewNode creates a new NodeStomp object and returns it.
func NewNode() *NodeStomp {

	n := &NodeStomp{transport.NewNode(),
		make(map[string]Subscription),
		make([]*HandlerStomp, 0)}

	return n
}

// SendFrame sends a frame to ActiveLink with certain ID.
func (n *NodeStomp) SendFrame(topic string, frame frame.Frame) {

	//log.WithField("topic", topic).Debugf("SendFrame()")
	log.Infof("SendFrame")

	n.CommandCh <- &CommandSendFrameStomp{
		transport.NodeCommand{Cmd: stompSendFrameCommand},
		frame,
		topic,
	}
}

/*
// GetChannel returns a channel connected with topic
func (n *NodeStomp) GetChannel(topic string) (chan frame.Frame, error) {
	var s Subscription
	var ok bool
	if s, ok = n.subscriptions[topic]; !ok {
		return nil, errors.New("No subscription for this topic")
	}
	return s.ch, nil
}
*/

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
