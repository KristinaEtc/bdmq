package stomp

import (
	"io"

	"github.com/KristinaEtc/bdmq/transport"
	"github.com/KristinaEtc/stomp/frame"
)

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
	link   transport.LinkWriter
	node   *transport.Node
	writer *frame.Writer
	reader *frame.Reader
}

// OnRead implements OnRead method from Heandler interface
func (h *HandlerStomp) OnRead(msg []byte) {

	// log.Debugf("OnRead msg=%s.", msg)
	// msgStr := strings.TrimSpace(string(msg))
	// buf := make([]byte, lA.linkDesc.bufSize)

	for {
		ff, err := h.reader.Read()
		if err != nil {
			if err == io.EOF {
				log.Errorf("connection closed: eof")
			} else {
				log.Errorf("read failed: %s", err.Error())
			}
			return
		}
		if ff == nil {
			log.Infof("heartbeat")
			continue
		}

		log.Debugf("Recieved data: [%s], [%v], [%s]", ff.Command, ff.Header, string(ff.Body))
		//log.Infof("message= %s", string(ff.Body))
	}

}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerStomp) OnConnect() error {

	log.Debugf("OnConnect %s %d", h.link.Id(), h.link.Mode())
	h.writer = frame.NewWriter(h.link.Conn())
	h.reader = frame.NewReader(h.link.Conn())

	return nil
}

// OnWrite implements OnWrite method from Heandler interface
func (h *HandlerStomp) OnWrite(msg []byte) {

	log.WithField("ID=", h.link.Id()).Debugf("OnWrite")

	var frameExample = frame.New(
		"SEND", "test1", "2", "message", string(msg))

	err := h.writer.Write(frameExample)
	if err != nil {
		log.Errorf("Err: %s", err.Error())
	}

	// Realize in disconnect()
	//f.Close()
}
