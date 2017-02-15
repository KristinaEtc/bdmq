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
func (h *HandlerEcho) OnRead(msg []byte) {

	//msgStr := strings.TrimSpace(string(msg))
	//log.Debugf("OnRead msg=%s. Resending it.", msgStr)
	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.Id()).Errorf("Error read: %s", err.Error())
	}

}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerEcho) OnConnect() error {

	log.Debugf("OnConnect %s %d", h.link.Id(), h.link.Mode())
	cnt := 65536
	tmp := make([]byte, cnt)
	for i := 0; i < cnt; i++ {
		tmp[i] = byte(48 + (i & 7))
	}
	tmp[cnt-1] = byte('\n')
	if h.link.Mode() == 1 {
		//h.OnWrite([]byte("0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF\n"))
		h.OnWrite(tmp)
	}
	return nil
}

// OnWrite implements OnWrite method from Heandler interface
func (h *HandlerEcho) OnWrite(msg []byte) {

	log.WithField("ID=", h.link.Id()).Debugf("OnWrite")
	err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.Id()).Errorf("Error read: %s", err.Error())
	}
}
