package stomp

import (
	"errors"
	"io"
	"time"

	"fmt"

	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

// DefaultHeartBeatError is a default time span to add to read/write heart-beat timeouts
// to avoid premature disconnections due to network latency.
const DefaultHeartBeatError = 5 * time.Second

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
	config Config

	options *connOptions

	version      Version
	readTimeout  time.Duration
	writeTimeout time.Duration
	server       string
	session      string
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

	switch frame.Command {
	case "CONNECT":
		{
			h.log.Infof("[%s]/[%s]: not implemented", frame.Command, frame.Header)
		}

	case "SEND":
		{
			h.log.Infof("[%s]/[%s]", frame.Command, frame.Header)
			h.node.CommandCh <- &CommandReceiveFrameStomp{
				transport.NodeCommand{Cmd: stompReceiveFrameCommand},
				frame,
				linkActiveID,
			}
		}
	case "SUBSCRIBE":
		{
			h.log.Infof("[%s]/[%s]: not implemented", frame.Command, frame.Header)
		}
	case "UNSUBSCRIBE":
		{
			h.log.Infof("[%s]/[%s]: not implemented", frame.Command, frame.Header)
		}
	case "BEGIN":
		{
			h.log.Infof("[%s]/[%s]: not implemented", frame.Command, frame.Header)
		}
	case "COMMIT":
		{
			h.log.Infof("[%s]/[%s]: not implemented", frame.Command, frame.Header)
		}
	case "ABORT":
		{
			h.log.Infof("[%s]/[%s]: not implemented", frame.Command, frame.Header)
		}
	case "ACK":
		{
			h.log.Infof("[%s]/[%s]: not implemented", frame.Command, frame.Header)
		}
	case "NACK":
		{
			h.log.Infof("[%s]/[%s]: not implemented", frame.Command, frame.Header)
		}
	case "DISCONNECT":
		{
			h.log.Infof("[%s]/[%s]: not implemented", frame.Command, frame.Header)
		}
	}
}

// OnRead implements OnRead method from transport.Heandler interface
func (h *HandlerStomp) OnRead(rd io.Reader) error {

	h.log.Info("HANDLER ONREAD")
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

func (h *HandlerStomp) readFrame(frameCommand string, rd io.Reader) (*frame.Frame, error) {
	//recieving//h.log.Infof("message= %v", string(msg))
	reader := frame.NewReader(rd)

	fr, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			h.log.Errorf("connection closed: eof")
		} else {
			h.log.Errorf("read failed: %s", err.Error())
		}
		return nil, err
	}
	if fr.Command != frameCommand {
		log.Errorf("Got frame with command [%s], expect [CONNECTED]", fr.Command)
		return nil, fmt.Errorf("Got frame with command [%s], expect [CONNECTED]", fr.Command)
	}

	return fr, nil
}

func (h *HandlerStomp) sendConnectedFrame() {
	//sending
}

func (h *HandlerStomp) recieveConnectedFrame() {
	//recieving
}

// OnConnect implements OnConnect method from Heandler interface
func (h *HandlerStomp) OnConnect(rd io.Reader) error {
	h.log.Debugf("OnConnect %s", h.link.Mode())

	var opts []func(h *HandlerStomp) error

	switch h.link.Mode() {
	case "client":
		{
			log.Warn("CLIENT tryin CONNECT")
			options, err := newConnOptions(h, opts)
			if err != nil {
				return err
			}

			if options.Host == "" {
				options.Host = "default"
			}

			connectFrame, err := options.NewFrame()
			if err != nil {
				return err
			}

			err = h.Writer.Write(connectFrame)
			if err != nil {
				return err
			}

			response, err := h.readFrame("CONNECTED", rd)
			if err != nil {
				return err
			}

			h.server = response.Header.Get(frame.Server)
			h.session = response.Header.Get(frame.Session)
			log.Debugf("h.server=[%s], h.session=[%s]", h.server, h.session)

			if versionString := response.Header.Get(frame.Version); versionString != "" {
				log.Debugf("versionString=[%s]", versionString)
				version := Version(versionString)
				if err = version.CheckSupported(); err != nil {
					return fmt.Errorf("Wrong version in CONNECTED frame: %s", err.Error())
				}
				h.version = version
				log.Debugf("h.version=[%s]", h.version)
			} else {
				// no version in the response, so assume version 1.0
				h.version = v10

				//TODO: add heartbeat
			}

		}
	case "server":
		{
			log.Warn("SERVER tryin CONNECTED")
			f, err := h.readFrame("CONNECT", rd)
			if err != nil {
				return err
			}
			log.Debugf("SERVER CONNECTED=[%v]/[%v]/[%v]", f, f.Header, f.Body)

			if _, ok := f.Header.Contains(frame.Receipt); ok {
				// CONNNECT and STOMP frames are not allowed to have
				// a receipt header.
				log.Errorf(" CONNNECT and STOMP frames are not allowed to have  a receipt header")
				return errors.New("CONNNECT and STOMP frames are not allowed to have a receipt header")
			}

			// if either of these fields are absent, pass nil to the
			// authenticator function.
			login, _ := f.Header.Contains(frame.Login)
			log.Debugf("login=[%s]", login)
			passcode, _ := f.Header.Contains(frame.Passcode)
			log.Debugf("passcode=[%s]", passcode)

			h.version, err = determineVersion(f)
			if err != nil {
				log.Error("protocol version negotiation failed")
				return err
			}

			log.Debugf("Version=[%s]", h.version)

			if h.version == v10 {
				// don't want to handle V1.0 at the moment
				// TODO: get working for V1.0
				log.Errorf("unsupported version %s", h.version)
				return fmt.Errorf("unsupported version %s", h.version)
			}

			cx, cy, err := getHeartBeat(f)
			if err != nil {
				log.Errorf("invalid heart-beat")
				return err
			}

			log.Debugf("heartBeat=[%d]/[%d]", cx/int(time.Millisecond), cy/int(time.Millisecond))

			response := frame.New(frame.CONNECTED,
				frame.Version, string(h.version),
				frame.Server, "stompd/x.y.z", // TODO: get version
				frame.HeartBeat, fmt.Sprintf("%d,%d", cy, cx))

			err = h.Writer.Write(response)
			if err != nil {
				return err
			}
		}
	}
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
