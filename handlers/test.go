package handlers

import (
	"bufio"
	"io"
	"strings"

	"github.com/KristinaEtc/bdmq/transport"
)

// HandlerTestFactory is a factory of Ha
type HandlerTestFactory struct {
}

// InitHandler creates a new HandlerTest and returns it.
// InitHandler is a function of transport.HandlerFactory.
func (h HandlerTestFactory) InitHandler(n *transport.Node, l transport.LinkWriter) transport.Handler {

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
func (h *HandlerTest) OnRead(rd io.Reader) error {

	var msg []byte
	var err error
	for {
		msg, err = bufio.NewReader(rd).ReadBytes('\n')
		if err != nil {
			log.Errorf("Error read: %s", err.Error())
			return err
		}
		//msgStr := strings.TrimSpace(string(message))
		//lA.log.Debugf("Message Received: %s", msgStr)

		msgStr := strings.TrimSpace(string(msg))
		log.WithField("ID=", h.link.ID()).Debugf("OnRead msg=%s", msgStr)
	}
}

// OnConnect implements OnConnect method from transport.Handler interface
func (h *HandlerTest) OnConnect(rd io.Reader) error {
	log.WithField("ID=", h.link.ID()).Debugf("OnConnect")
	return nil
}

// OnWrite implements OnWrote method from Handler interface
func (h *HandlerTest) OnWrite(msg []byte) {

	log.WithField("ID=", h.link.ID()).Debugf("OnWrite")

	_, err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.ID()).Errorf("Error read: %s", err.Error())
	}
}

// OnDisconnect implements OnDisconnect method from Handler interface
func (h *HandlerTest) OnDisconnect() {
	log.WithField("ID=", h.link.ID()).Debugf("OnConnect")
}
