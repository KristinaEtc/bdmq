package stomp

import (
	"errors"
	"io"

	"github.com/KristinaEtc/bdmq/transport"
	"github.com/go-stomp/stomp/frame"
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
		link:     l,
		node:     n,
		Writer:   frame.NewWriter(wd),
		Reader:   frame.NewReader(rd),
		log:      slf.WithContext("stompHandler").WithFields(slf.Fields{"ID": l.ID()}),
		topicChs: make(map[string]*chan transport.Frame),
	}
	return handler
}

// HandlerStomp realize Handler interface from transport package
type HandlerStomp struct {
	link     transport.LinkWriter
	node     *transport.Node
	Reader   *frame.Reader
	Writer   *frame.Writer
	log      slf.Logger
	topicChs map[string]*chan transport.Frame
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
func (h *HandlerStomp) OnRead() {

	//h.log.Infof("message= %v", string(msg))

	for {
		frame, err := h.Reader.Read()
		if err != nil {
			if err == io.EOF {
				h.log.Errorf("connection closed: eof")
			} else {
				h.log.Errorf("read failed: %s", err.Error())
			}
			return
		}
		if frame == nil {
			h.log.Infof("heartbeat")
			continue
		}

		processFrame("test", frame)
	}

}

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

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerStomp) OnConnect() error {
	h.log.Debugf("OnConnect %d", h.link.Mode())
	return nil
}

// OnWrite implements OnWrite method from transport.Heandler interface
func (h *HandlerStomp) OnWrite(frame *transport.Frame) {

	h.log.Debug("stomp handler Write")
	h.Writer.Write(frame)
	return
}

// OnDisconnect implements OnDisconnect method from transport.Heandler interface
func (h *HandlerStomp) OnDisconnect() {

	h.log.Debugf("OnDisconnect %d", h.link.Mode())
	//f.Close()
}
