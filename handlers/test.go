package handlers

import "github.com/KristinaEtc/bdmq/transport"

type HandlerTestFactory struct {
}

func (h HandlerTestFactory) InitHandler(l *transport.LinkActive, n *transport.Node) transport.Handler {

	log.Debugf("InitHandler")
	handler := &HandlerTest{
		link: l,
		node: n,
	}
	return handler
}

// HandlerTest realize Handler interface from transport package
type HandlerTest struct {
	link *transport.LinkActive
	node *transport.Node
}

// OnRead implements OnRead method from Heandler interface
func (h *HandlerTest) OnRead(msg string) {
	log.Debugf("OnRead msg=%s", msg)

}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerTest) OnConnect() error {
	log.Debugf("OnConnect")
	return nil
}

// OnWrite implements OnWrote method from Heandler interface
func (h *HandlerTest) OnWrite(msg string) {

	log.WithField("ID=", h.link.LinkActiveID).Debugf("OnWrite")

	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.LinkActiveID).Errorf("Error read: %s", err.Error())
	}
}
