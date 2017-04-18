package stomp

import (
	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
)

//CommandSendFrameStomp is a struct whith will send to StompNode proseccor
type CommandSendFrameStomp struct {
	transport.NodeCommand
	frame frame.Frame
	topic string
}

//CommandReceiveFrameStomp is a struct whith will send to StompNode proseccor
type CommandReceiveFrameStomp struct {
	transport.NodeCommand
	frame        frame.Frame
	linkActiveID string
}

// CommandRegisterHandlerStomp is a struct for StompNode processor for registration a new handler
type CommandRegisterHandlerStomp struct {
	transport.NodeCommand
	handler *HandlerStomp
}

//CommandSubscribeStomp is a struct whith will process all abut subscriptions
type CommandSubscribeStomp struct {
	transport.NodeCommand
	topic string
	sub   chan subscripted
}

const (
	stompSendFrameCommand            transport.CommandID = 100
	stompSubscribeCommand            transport.CommandID = 101
	stompUnsubscribeCommand          transport.CommandID = 102
	stompReceiveFrameCommand         transport.CommandID = 103
	stompRegisterStompHandlerCommand transport.CommandID = 104
)
