package stomp

import "github.com/KristinaEtc/bdmq/transport"

type HandlerStompFactory struct {
}

var StompFactory *HandlerStompFactory

func (h HandlerStompFactory) InitHandler(l transport.LinkWriter, n *transport.Node) transport.Handler {

	log.Debugf("InitHandler")
	handler := &HandlerStomp{
		link: l,
		node: n,
	}
	return handler
}

// HandlerEcho realize Handler interface from transport package
type HandlerStomp struct {
	link transport.LinkWriter
	node *transport.Node
}

// OnRead implements OnRead method from Heandler interface
func (h *HandlerStomp) OnRead(msg []byte) {
	log.Infof("message= %v", string(msg))
}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerStomp) OnConnect() error {
	log.Debugf("OnConnect %s %d", h.link.Id(), h.link.Mode())
	return nil
}

// OnWrite implements OnWrite method from Heandler interface
func (h *HandlerStomp) OnWrite(msg []byte) {

	log.WithField("ID=", h.link.Id()).Debugf("OnWrite")
	h.link.Write(msg)

	// Realize in disconnect()
	//f.Close()
}
