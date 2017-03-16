package handlers

import (
	"strings"
	"time"

	"github.com/KristinaEtc/bdmq/transport"
)

// HandlerHelloWorldFactory is a factory of HandlerHelloWorld
type HandlerHelloWorldFactory struct {
}

// InitHandler creates a new HandlerHelloWorld and returns it.
// InitHandler is a function of transport.HandlerFactory.
func (h HandlerHelloWorldFactory) InitHandler(l transport.LinkWriter, n *transport.Node) transport.Handler {

	log.Debugf("InitHandler")
	handler := &HandlerHelloWorld{
		link: l,
		node: n,
	}
	return handler
}

// HandlerHelloWorld realize Handler interface from transport package
type HandlerHelloWorld struct {
	link transport.LinkWriter
	node *transport.Node
}

// OnRead implements OnRead method from transport.Handler interface
func (h *HandlerHelloWorld) OnRead(msg []byte) {

	msgStr := strings.TrimSpace(string(msg))
	log.WithField("ID=", h.link.ID()).Debugf("OnRead msg=%s. Resending it.", msgStr)
}

// OnConnect implements OnConnect method from transport.Handler interface
func (h *HandlerHelloWorld) OnConnect() error {
	log.WithField("ID=", h.link.ID()).Debugf("OnConnect")
	// TODO: to add channel where will be sended a kill signal
	ticker := time.NewTicker(time.Second * 2)
	func() {
		for _ = range ticker.C {
			h.OnWrite([]byte("Hello World\n"))
		}
	}()
	return nil
}

// OnWrite implements OnWrote method from transport.Handler interface
func (h *HandlerHelloWorld) OnWrite(msg []byte) {

	log.WithField("ID=", h.link.ID()).Debugf("OnWrite")
	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.ID()).Errorf("Error read: %s", err.Error())

	}
}

// OnDisconnect implements OnDisconnect method from Handler interface
func (h *HandlerHelloWorld) OnDisconnect() {
	log.WithField("ID=", h.link.ID()).Debugf("OnConnect")
}
