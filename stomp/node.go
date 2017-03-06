package stomp

import "github.com/KristinaEtc/bdmq/transport"

func NewNode() (n *transport.Node) {
	n = transport.NewNode()

	StompProcesser := &StompProcesser{node: n}
	n.AddCmdProcessor(StompProcesser)

	return n
}
