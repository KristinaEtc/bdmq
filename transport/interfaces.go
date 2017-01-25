package transport

import "github.com/ventu-io/slf"

var log = slf.WithContext("transport.go")

// Linker is an interface which present connection node and
// it's options.
type Linker interface {
	Write(b []byte) (int, error)
	//Disconnect()
	//GetStatus() string
}

// Handler is an interface for loop coupling between transport layer and
// protocol's realization.
type Handler interface {
	OnRead()
	OnConnect() error
	OnWrite()
	//...
}

type HandlerFactory interface {
	InitHandler(*LinkActive, *Node) Handler
}

// HandlerFactories stores loaded Handlers
type HandlerFactories map[string]HandlerFactory

var handlers HandlerFactories = make(map[string]HandlerFactory)

// RegisterHandlerFactory set Handlers
func RegisterHandlerFactory(hName string, h HandlerFactory) {
	log.Debug("RegisterHandlerFactory")
	handlers[hName] = h
}
