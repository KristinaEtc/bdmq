package stomp

import "github.com/KristinaEtc/bdmq/transport"

// HandlerStompFactory is a factory of HandlerStomp
type HandlerStompFactory struct {
}

// InitHandler creates a new HandlerStomp and returns it.
// InitHandler is a function of transport.HandlerFactory.
func (h HandlerStompFactory) InitHandler(l transport.LinkWriter, n *transport.Node) transport.Handler {

	log.Debugf("InitHandler")
	handler := &HandlerStomp{
		link: l,
		node: n,
	}
	return handler
}

// HandlerStomp realize Handler interface from transport package
type HandlerStomp struct {
	link transport.LinkWriter
	node *transport.Node
}

// OnRead implements OnRead method from transport.Heandler interface
func (h *HandlerStomp) OnRead(msg []byte) {
	log.WithField("ID=", h.link.ID()).Infof("message= %v", string(msg))
}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerStomp) OnConnect() error {
	log.WithField("ID=", h.link.ID()).Debugf("OnConnect %d", h.link.Mode())
	return nil
}

// OnWrite implements OnWrite method from transport.Heandler interface
func (h *HandlerStomp) OnWrite(msg []byte) {

	log.WithField("ID=", h.link.ID()).Debugf("OnWrite")
	h.link.Write(msg)
}

// OnDisconnect implements OnDisconnect method from transport.Heandler interface
func (h *HandlerStomp) OnDisconnect() {

	log.WithField("ID=", h.link.ID()).Debugf("OnDisconnect %d", h.link.Mode())
	//f.Close()
}
