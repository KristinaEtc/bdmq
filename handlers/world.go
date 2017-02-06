package handlers

import (
	"time"

	"github.com/KristinaEtc/bdmq/transport"
)

type HandlerHelloWorldFactory struct {
}

func (h HandlerHelloWorldFactory) InitHandler(l transport.LinkWriter, n *transport.Node) transport.Handler {

	log.Debugf("InitHandler")
	handler := &HandlerHelloWorld{
		link: l,
		node: n,
	}
	return handler
}

// HandlerEcho realize Handler interface from transport package
type HandlerHelloWorld struct {
	link transport.LinkWriter
	node *transport.Node
}

// OnRead implements OnRead method from Heandler interface
func (h *HandlerHelloWorld) OnRead(msg string) {
	log.Debugf("OnRead msg=%s. Resending it.", msg)
}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerHelloWorld) OnConnect() error {
	log.Debugf("OnConnect")
	// TODO: to add channel where will be sended a kill signal
	ticker := time.NewTicker(time.Second * 2)
	func() {
		for _ = range ticker.C {
			h.OnWrite("Hello World\n")
		}
	}()
	return nil
}

// OnWrite implements OnWrote method from Heandler interface
func (h *HandlerHelloWorld) OnWrite(msg string) {

	log.WithField("ID=", h.link.GetActiveLinkID()).Debugf("OnWrite")
	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.GetActiveLinkID()).Errorf("Error read: %s", err.Error())

	}
}
