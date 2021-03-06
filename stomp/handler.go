package stomp

import (
	"io"

	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

// HandlerStompFactory is a factory of HandlerStomp
type HandlerStompFactory struct {
}

// InitHandler creates a new HandlerStomp and returns it.
// InitHandler is a function of transport.HandlerFactory.
func (h HandlerStompFactory) InitHandler(n *transport.Node, l transport.LinkWriter) transport.Handler {

	log.Debugf("HandlerStompFactory.InitHandler() enter")
	handler := &HandlerStomp{
		link:   l,
		Writer: frame.NewWriter(l),
		log:    slf.WithContext("stompHandler").WithFields(slf.Fields{"ID": l.ID()}),
		node:   n,
		//topicChs: make(map[string]*chan transport.Frame),
	}
	return handler
}

// HandlerStomp realize Handler interface from transport package
type HandlerStomp struct {
	link   transport.LinkWriter
	node   *transport.Node
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

// receiveFrame used for recieving a frame to stompProcessor.
func (h *HandlerStomp) receiveFrame(linkActiveID string, frame frame.Frame) {
	//h.log.Debugf("func ReceiveFrame()")

	h.node.CommandCh <- &CommandReceiveFrameStomp{
		transport.NodeCommand{Cmd: stompReceiveFrameCommand},
		frame,
		linkActiveID,
	}
}

// OnRead implements OnRead method from transport.Heandler interface
func (h *HandlerStomp) OnRead(rd io.Reader) error {

	//h.log.Infof("message= %v", string(msg))
	reader := frame.NewReader(rd)

	for {
		fr, err := reader.Read()
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
		h.receiveFrame(h.link.ID(), *fr)
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

// receiveFrame used for recieving a frame to stompProcessor.
func (h *HandlerStomp) registerStompHandler() {
	h.log.Debugf("func registerStompHandler()")

	h.node.CommandCh <- &CommandRegisterHandlerStomp{
		transport.NodeCommand{Cmd: stompRegisterStompHandlerCommand},
		h,
	}
}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerStomp) OnConnect() error {
	h.log.Debugf("OnConnect %d", h.link.Mode())
	h.registerStompHandler()
	return nil
}

// OnWrite implements OnWrite method from transport.Heandler interface
func (h *HandlerStomp) OnWrite(frame frame.Frame) {

	//h.log.Debug("OnWrite")
	h.Writer.Write(&frame)
	//h.log.Debug("OnWrite exit")
	return
}

// OnDisconnect implements OnDisconnect method from transport.Heandler interface
func (h *HandlerStomp) OnDisconnect() {

	h.log.Debugf("OnDisconnect %d", h.link.Mode())
	//f.Close()
}
