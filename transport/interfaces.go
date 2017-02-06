package transport

import (
	"math/rand"
	"time"

	"github.com/ventu-io/slf"
)

var log = slf.WithContext("transport.go")

// SizeOfBuf is a size of tcp buffer
const SizeOfBuf = 1024

const (
	commandQuit = "QUIT"
)

//for generating reconnect endpoint
var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

/*
// Linker is an interface which present connection node and
// it's options.
type Linker interface {
	Write(string) error
	Read()
	//Disconnect()
	//GetStatus() string
}
*/

// Handler is an interface for loop coupling between transport layer and
// protocol's realization.
type Handler interface {
	OnRead(msg string)
	OnConnect() error
	OnWrite(msg string)
	//...
}

type HandlerFactory interface {
	InitHandler(LinkActiver, *Node) Handler
}

// HandlerFactories stores loaded Handlers
type HandlerFactories map[string]HandlerFactory

var handlers HandlerFactories = make(map[string]HandlerFactory)

// RegisterHandlerFactory set Handlers
func RegisterHandlerFactory(hName string, h HandlerFactory) {
	log.Debug("RegisterHandlerFactory")
	handlers[hName] = h
}
