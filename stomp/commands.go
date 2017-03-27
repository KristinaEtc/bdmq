package stomp

import "github.com/KristinaEtc/bdmq/transport"

//CommandSendFrameStomp is a struct with will send to StompNode proseccor
type CommandSendFrameStomp struct {
	transport.NodeCommand
	frame        transport.Frame
	linkActiveID string
}

//CommandSubscribeStomp is a struct with will process all abut subscriptions
type CommandSubscribeStomp struct {
	transport.NodeCommand
	topic        string
	linkActiveID string
	sub          chan subscripted
}

const (
	stompSendFrameCommand transport.CommandID = 100
	stompSubscribeCommand transport.CommandID = 101
)
