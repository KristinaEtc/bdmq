package handlers

import "github.com/KristinaEtc/bdmq/transport"
import "github.com/ventu-io/slf"

var log = slf.WithContext("handlers.go")

type HandlerEchoFactory struct {
}

var EchoFactory *HandlerEchoFactory

func (h HandlerEchoFactory) InitHandler(l transport.LinkWriter, n *transport.Node) transport.Handler {

	log.Debugf("InitHandler")
	handler := &HandlerEcho{
		link: l,
		node: n,
	}
	return handler
}

// HandlerEcho realize Handler interface from transport package
type HandlerEcho struct {
	link transport.LinkWriter
	node *transport.Node
}

// OnRead implements OnRead method from Heandler interface
func (h *HandlerEcho) OnRead(msg string) {
	log.Debugf("OnRead msg=%s. Resending it.", msg)
	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.GetActiveLinkID()).Errorf("Error read: %s", err.Error())
	}

}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerEcho) OnConnect() error {
	log.Debugf("OnConnect")
	return nil
}

// OnWrite implements OnWrote method from Heandler interface
func (h *HandlerEcho) OnWrite(msg string) {

	log.WithField("ID=", h.link.GetActiveLinkID()).Debugf("OnWrite")

	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.GetActiveLinkID()).Errorf("Error read: %s", err.Error())
	}
}
