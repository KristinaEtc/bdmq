package stomp

import (
	"io"

	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

/*
type Frame interface{
	Dump()
}
*/

// HandlerStompFactory is a factory of HandlerStomp
type HandlerStompFactory struct {
}

// InitHandler creates a new HandlerStomp and returns it.
// InitHandler is a function of transport.HandlerFactory.
func (h HandlerStompFactory) InitHandler(l transport.LinkWriter, n *transport.Node, rd io.Reader, wd io.Writer) transport.Handler {

	log.Debugf("HandlerStompFactory.InitHandler() enter")
	handler := &HandlerStomp{
		link: l,
		//node:   n,
		Writer: frame.NewWriter(wd),
		Reader: frame.NewReader(rd),
		log:    slf.WithContext("stompHandler").WithFields(slf.Fields{"ID": l.ID()}),
		//topicChs: make(map[string]*chan transport.Frame),
	}
	return handler
}

// HandlerStomp realize Handler interface from transport package
type HandlerStomp struct {
	link   transport.LinkWriter
	node   NodeStomp
	Reader *frame.Reader
	Writer *frame.Writer
	log    slf.Logger
	//	topicChs map[string]*chan transport.Frame
}

// processFrame add frame for indicated topic
func processFrame(topic string, frame *frame.Frame) error {
	/*
		for _, header := range *frame.Header.GetAll(topic) {
			if strings.Compare(header, "topic:"+topic) == 0 {

			}
		}
		topicName, ok := frame.Header[topic]
		if !ok {
			return error.New("Got message without topic header; ignored")
		}
		// TODO: parse several topics by commas
		p.topicChs[topicName] <- []byte(frame)
	*/
	return nil
}

// OnRead implements OnRead method from transport.Heandler interface
func (h *HandlerStomp) OnRead() error {

	//h.log.Infof("message= %v", string(msg))

	for {
		fr, err := h.Reader.Read()
		if err != nil {
			if err == io.EOF {
				h.log.Errorf("connection closed: eof")
			} else {
				h.log.Errorf("read failed: %s", err.Error())
			}
			return err
		}
		if fr == nil {
			h.log.Infof("heartbeat")
			continue
		}

		//	h.log.Infof("message= %s", fr.Dump())
		h.node.ReceiveFrame(h.link.ID(), "test-topic", *fr)
	}

}

/*
// Subscribe returns channel for topics receiving
func (h *HandlerStomp) Subscribe(topic string) (*chan transport.Frame, error) {
	h.log.Debugf("Subscribe")

	_, ok := h.topicChs[topic]
	if ok {
		log.Warnf("Channel with such topic have already exists.")
		return nil, errors.New("Channel with such topic have already exists.")
	}

	ch := make(chan transport.Frame, 0)
	h.topicChs[topic] = &ch
	log.Debug("Subscribed successfully")
	return &ch, nil
}
*/

func (h *HandlerStomp) initNodeStomp() {
	h.node = *nodes["NodeStomp-ID"]
}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerStomp) OnConnect() error {
	h.log.Debugf("OnConnect %d", h.link.Mode())
	h.initNodeStomp()
	return nil
}

// OnWrite implements OnWrite method from transport.Heandler interface
func (h *HandlerStomp) OnWrite(frame frame.Frame) {

	h.log.Debug("OnWrite")
	h.Writer.Write(&frame)
	h.log.Debug("OnWrite exit")
	return
}

// OnDisconnect implements OnDisconnect method from transport.Heandler interface
func (h *HandlerStomp) OnDisconnect() {

	h.log.Debugf("OnDisconnect %d", h.link.Mode())
	//f.Close()
}
