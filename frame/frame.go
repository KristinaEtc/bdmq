package frame

import (
	"bytes"
	"fmt"
)

/*

// FactoryStomp is a factory for creating ProcessorStomp
type FactoryStomp struct {
}

// ProcessorStomp is an interface for FrameProcessors
type ProcessorStomp struct {
	linkActive transport.LinkActive
	Reader     *Reader
	Writer     *Writer
	log        slf.Logger
	topicCh    *chan []byte
}

// InitFrameProcessor creates a new entity of ProcessorStomp and adds it to Node process slice
func (f FactoryStomp) InitFrameProcessor(lActive transport.LinkActive, rd io.Reader, wd io.Writer) transport.FrameProcessor {

	sFrameProcessor := &ProcessorStomp{
		Writer:     NewWriter(wd),
		Reader:     NewReader(rd),
		log:        slf.WithContext("FrameProcessor").WithFields(slf.Fields{"ID": lActive.ID()}),
		linkActive: lActive,
		topicCh:    lActive.GetTopicCh(),
	}

	return sFrameProcessor
}

func (p *ProcessorStomp) write(data interface{}) {
	p.log.Debug("stomp Processor Write")
	p.Writer.Write(data.(*Frame))
	return
}

// ToByte converts Frame (for different FrameProcessers different types) type to slices of bytes.
// Implements ToByte method from transport.FrameProcessor interface.
func (p *ProcessorStomp) ToByte(data interface{}) []byte {
	p.log.Debug("ToByte()")
	/*b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	f := data.(Frame)
	p.Writer.Write(&f)
	return []byte("data.(*Frame)")
}

func (p *ProcessorStomp) Read() error {
	for {
		ff, err := p.Reader.Read()
		if err != nil {
			if err == io.EOF {
				p.log.Errorf("connection closed: eof")
			} else {
				p.log.Errorf("read failed: %s", err.Error())
			}
			return err
		}
		if ff == nil {
			p.log.Infof("heartbeat")
			continue
		}

		topicName, ok := ff.Header[topic]
		if !ok {
			log.Warn("Got message without topic header; ignored")
			continue
		}
		// TODO: parse by commas
		p.topicChs[topicName] <- []byte(ff.Dump())

	}
}

*/

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

// Dump prints all frame content.
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
