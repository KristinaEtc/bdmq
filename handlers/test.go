package handlers

import (
	"strings"

	"github.com/KristinaEtc/bdmq/transport"
)

// HandlerTestFactory is a factory of Ha
type HandlerTestFactory struct {
}

// InitHandler creates a new HandlerTest and returns it.
// InitHandler is a function of transport.HandlerFactory.
func (h HandlerTestFactory) InitHandler(l transport.LinkWriter, n *transport.Node) transport.Handler {

	log.Debugf("InitHandler")
	handler := &HandlerTest{
		link: l,
		node: n,
	}
	return handler
}

// HandlerTest realize Handler interface from transport package
type HandlerTest struct {
	link transport.LinkWriter
	node *transport.Node
}

// OnRead implements OnRead method from Heandler interface
func (h *HandlerTest) OnRead(msg []byte) {

	msgStr := strings.TrimSpace(string(msg))
	log.WithField("ID=", h.link.ID()).Debugf("OnRead msg=%s", msgStr)
}

// OnConnect implements OnConnect method from transport.Handler interface
func (h *HandlerTest) OnConnect() error {
	log.WithField("ID=", h.link.ID()).Debugf("OnConnect")
	return nil
}

// OnWrite implements OnWrote method from Handler interface
func (h *HandlerTest) OnWrite(msg []byte) {

	log.WithField("ID=", h.link.ID()).Debugf("OnWrite")

	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.ID()).Errorf("Error read: %s", err.Error())
	}
}

// OnDisconnect implements OnDisconnect method from Handler interface
func (h *HandlerTest) OnDisconnect() {
	log.WithField("ID=", h.link.ID()).Debugf("OnConnect")
}
