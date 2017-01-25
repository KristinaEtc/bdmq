package handler

import "github.com/KristinaEtc/bdmq/transport"
import "github.com/ventu-io/slf"

var log = slf.WithContext("handler-test.go")

type HandlerTestFactory struct {
}

var TestFactory *HandlerTestFactory

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
func (h *HandlerTest) OnRead() {
	log.Debugf("OnRead")

}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerTest) OnConnect() error {
	log.Debugf("OnConnect")
	return nil
}

// OnWrite implements OnWrote method from Heandler interface
func (h *HandlerTest) OnWrite() {
	log.Debugf("OnWrite")
	//return nil
}
