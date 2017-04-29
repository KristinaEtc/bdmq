package handlers

import (
	"bufio"
	"io"

	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("test handlers")

// HandlerEchoFactory is a factory of HandlerEcho
type HandlerEchoFactory struct {
}

// InitHandler creates a new HandlerEcho and returns it.
// InitHandler is a function of transport.HandlerFactory.
func (h HandlerEchoFactory) InitHandler(n *transport.Node, l transport.LinkWriter) transport.Handler {

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

// OnRead implements OnRead method from transport.Handler interface
func (h *HandlerEcho) OnRead(rd io.Reader) error {

	var msg []byte
	var err error
	for {
		msg, err = bufio.NewReader(rd).ReadBytes('\n')
		if err != nil {
			log.Errorf("Error read: %s", err.Error())
			return err
		}
		//msgStr := strings.TrimSpace(string(msg))
		//log.Debugf("OnRead msg=%s. Resending it.", msgStr)
		_, err = h.link.Write(msg)
		if err != nil {
			log.WithField("ID=", h.link.ID()).Errorf("Error read: %s", err.Error())
			return err
		}
	}
}

// OnConnect implements OnConnect method from transport.Handler interface
func (h *HandlerEcho) OnConnect(rd io.Reader) error {

	log.WithField("ID=", h.link.ID()).Debugf("OnConnect %s", h.link.Mode())

	cnt := 65536
	tmp := make([]byte, cnt)
	for i := 0; i < cnt; i++ {
		tmp[i] = byte(48 + (i & 7))
	}
	tmp[cnt-1] = byte('\n')
	if h.link.Mode() == "client" {
		//h.OnWrite([]byte("0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF\n"))
		h.OnWrite(tmp)
	}
	return nil
}

// OnWrite implements OnWrite method from transport.Handler interface
func (h *HandlerEcho) OnWrite(msg []byte) {

	log.WithField("ID=", h.link.ID()).Debugf("OnWrite")
	_, err := h.link.Write(msg)
	if err != nil {
		log.WithField("ID=", h.link.ID()).Errorf("Error read: %s", err.Error())
	}
}

// OnDisconnect implements OnDisconnect method from Handler interface
func (h *HandlerEcho) OnDisconnect() {
	log.WithField("ID=", h.link.ID()).Debugf("OnConnect")
}
