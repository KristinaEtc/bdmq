package transport

import (
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("transport.go")

type HandlerFactories map[string]HandlerFactory

var handlers HandlerFactories = make(map[string]HandlerFactory)

type Handler interface {
	OnRead([]byte)
	OnConnect() error
	OnWrite([]byte)
}

type HandlerFactory interface {
	InitHandler(LinkWriter, *Node) Handler
}

func RegisterHandlerFactory(handlerName string, hFactory HandlerFactory) {
	log.Debug("func RegisterHandlerFactory()")
	handlers[handlerName] = hFactory
}
