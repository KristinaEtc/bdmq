package handlers

import "github.com/KristinaEtc/bdmq/transport"
import "github.com/ventu-io/slf"
import "strings"

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
func (h *HandlerEcho) OnRead(msg []byte) {

	msgStr := strings.TrimSpace(string(msg))
	log.Debugf("OnRead msg=%s. Resending it.", msgStr)
	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.Id()).Errorf("Error read: %s", err.Error())
	}

}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerEcho) OnConnect() error {
	log.Debugf("OnConnect")
	return nil
}

// OnWrite implements OnWrote method from Heandler interface
func (h *HandlerEcho) OnWrite(msg []byte) {

	log.WithField("ID=", h.link.Id()).Debugf("OnWrite")
	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.Id()).Errorf("Error read: %s", err.Error())
	}
}
