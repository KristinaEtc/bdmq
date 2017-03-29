package transport

import (
	"io"

	"github.com/ventu-io/slf"
)

var log = slf.WithContext("transport")

// handlerFactories
type handlerFactories map[string]HandlerFactory

var handlers handlerFactories = make(map[string]HandlerFactory)

// Handler is an interface for handlers
// ...
type Handler interface {
	OnRead(io.Reader) error // read method
	OnConnect() error       // connect method
	//OnWrite([]byte)   // write method
	OnDisconnect() // disconnect method
	//Subscribe(string) (*chan Frame, error)
}

// HandlerFactory is an interface for creating new Handler
type HandlerFactory interface {
	//InitHandler(LinkWriter, *Node) Handler // creates new Handler
	InitHandler(*Node, LinkWriter) Handler
}

// RegisterHandlerFactory added HandlerFactory hFactory with name handlerName.
// Further this Handler could be used with this LinkActives
func RegisterHandlerFactory(handlerName string, hFactory HandlerFactory) {
	log.Debug("func RegisterHandlerFactory()")
	handlers[handlerName] = hFactory
}
