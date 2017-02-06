package handlers

import "github.com/KristinaEtc/bdmq/transport"

type HandlerTestFactory struct {
}

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
func (h *HandlerTest) OnRead(msg string) {
	log.Debugf("OnRead msg=%s", msg)
	/*
		if msg==commandQuit{
			log.Debug("got command quit for ActiveLink with ID=%s ; exiting")
		}
	*/

}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerTest) OnConnect() error {
	log.Debugf("OnConnect")
	return nil
}

// OnWrite implements OnWrote method from Heandler interface
func (h *HandlerTest) OnWrite(msg string) {

	log.WithField("ID=", h.link.GetActiveLinkID()).Debugf("OnWrite")

	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.GetActiveLinkID()).Errorf("Error read: %s", err.Error())
	}
}
