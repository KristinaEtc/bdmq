/*
Package frame provides functionality for manipulating STOMP frames.
*/
package frame

import (
	"bytes"
	"fmt"
	"io"

	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

type StompFrameProcessorFactory struct {
}

var sFrameProcessorFactory StompFrameProcessorFactory

type StompFrameProcessor struct {
	linkActive transport.LinkActive
	Reader     *Reader
	Writer     *Writer
	log        slf.Logger
}

func (s StompFrameProcessorFactory) InitFrameProcessor(lActive transport.LinkActive, rd io.Reader, wd io.Writer, log slf.Logger) transport.FrameProcessor {

	log.Debugf("InitStompFrame")
	sFrameProcessor := &StompFrameProcessor{
		Writer:     NewWriter(wd),
		Reader:     NewReader(rd),
		log:        log,
		linkActive: lActive,
	}

	return sFrameProcessor
}

func (s *StompFrameProcessor) Write(data interface{}) {
	s.log.Debug("stomp Processor Write")
	s.Writer.Write(data.(*Frame))
	return
}

func (s *StompFrameProcessor) ToByte(data interface{}) []byte {
	s.log.Debug("stomp To Byte")
	/*b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}*/
	f := data.(Frame)
	s.Writer.Write(&f)
	return []byte("data.(*Frame)")
}

func (s *StompFrameProcessor) Read() error {
	for {
		ff, err := s.Reader.Read()
		if err != nil {
			if err == io.EOF {
				s.log.Errorf("connection closed: eof")
			} else {
				s.log.Errorf("read failed: %s", err.Error())
			}
			return err
		}
		if ff == nil {
			s.log.Infof("heartbeat")
			continue
		}
		// TODO: change it to a channel
		s.log.Infof("2 Data: [%s], [%v], [%s]", ff.Command, ff.Header, string(ff.Body))
	}
}

// A Frame represents a STOMP frame. A frame consists of a command
// followed by a collection of header entries, and then an optional
// body.
type Frame struct {
	Command string
	Header  *Header
	Body    []byte
}

// New creates a new STOMP frame with the specified command and headers.
// The headers should contain an even number of entries. Each even index is
// the header name, and the odd indexes are the assocated header values.
func New(command string, headers ...string) *Frame {
	f := &Frame{Command: command, Header: &Header{}}
	for index := 0; index < len(headers); index += 2 {
		f.Header.Add(headers[index], headers[index+1])
	}
	return f
}

// Clone creates a deep copy of the frame and its header. The cloned
// frame shares the body with the original frame.
func (f *Frame) Clone() *Frame {
	fc := &Frame{Command: f.Command}
	if f.Header != nil {
		fc.Header = f.Header.Clone()
	}
	if f.Body != nil {
		fc.Body = make([]byte, len(f.Body))
		copy(fc.Body, f.Body)
	}
	return fc
}

func (f *Frame) String() string {
	return fmt.Sprintf("%s %s", f.Command, string(f.Body))
}

func (f *Frame) Dump() string {
	var buffer bytes.Buffer
	for index := 0; index < f.Header.Len(); index++ {
		if index > 0 {
			buffer.WriteString(", ")
		}
		k, v := f.Header.GetAt(index)
		fmt.Fprintf(&buffer, "\"%s\": \"%s\"", k, v)
		//buffer.WriteString(fmt.Sprintf(
		//f.Header.Add(headers[index], headers[index+1])
	}
	return fmt.Sprintf("%s %s %s", f.Command, string(f.Body), buffer.String())
}
