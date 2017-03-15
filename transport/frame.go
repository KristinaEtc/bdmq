package transport

import (
	"bufio"
	"io"
	"net"

	"github.com/ventu-io/slf"
)

type FrameProcessorFactories map[string]FrameProcessorFactory
type Frame interface{}

var frameProcessors FrameProcessorFactories = make(map[string]FrameProcessorFactory)

type FrameProcessor interface {
	//NewReader(rd io.Reader) io.Reader
	//	NewWriter() io.Writer
	Read() error
	//Write(interface{})
	ToByte(interface{}) []byte
}

type FrameProcessorFactory interface {
	InitFrameProcessor(LinkActive, io.Reader, io.Writer, slf.Logger) FrameProcessor
}

func RegisterFrameProcessorFactory(frameProcesserName string, fFactory FrameProcessorFactory) {
	log.Debug("func RegisterFrameFactory()")
	frameProcessors[frameProcesserName] = fFactory
}

//-------------------------------------------------------

type DefaultFrameProcessorFactory struct {
}

var dFrameProcessorFactory DefaultFrameProcessorFactory

// DefaultFrameProcesser implements FrameProcesser interface.
// Used by default.
type DefaultFrameProcessor struct {
	linkActive LinkActive
	reader     bufio.Reader
	writer     io.Writer
	log        slf.Logger
	handler    Handler
}

func (h DefaultFrameProcessorFactory) InitFrameProcessor(lActive LinkActive, rd net.Conn, wd net.Conn, log slf.Logger) FrameProcessor {

	dFrameProcesser := &DefaultFrameProcessor{
		writer:  wd,
		log:     log,
		handler: lActive.getHandler(),
	}

	dFrameProcesser.NewReader(rd)

	return dFrameProcesser
}

func (d *DefaultFrameProcessor) NewReader(rd io.Reader) (r io.Reader) {
	r = bufio.NewReader(rd)
	return
}

//func (d *DefaultFrameProcessor) Read() (message []byte, err error) {
func (d *DefaultFrameProcessor) Read() error {
	for {
		message, err := d.reader.ReadBytes('\n')
		if err != nil {
			d.log.Errorf("Error read: %s", err.Error())
			return err
		}
		d.handler.OnRead(message)
	}
}

func (d *DefaultFrameProcessor) ToByte(msg interface{}) []byte {
	return (msg.([]byte))
}
