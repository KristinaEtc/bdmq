package frame

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ventu-io/slf"
)

var pwdCurr = "frame"
var log = slf.WithContext(pwdCurr)

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

	var len = len(headers)
	/*	if len%2 != 0 {
			log.Errorf("Odd number of headers: set [key, value] headers; last with be ignored")
			len--
		}
	*/
	for index := 0; index < len; index += 2 {
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

// CompareFrames compares all values if frames
func CompareFrames(fr1, fr2 Frame) bool {

	if strings.Compare(fr1.Command, fr2.Command) == 0 {
		if bytes.Compare(fr1.Body, fr2.Body) == 0 {
			lenH := len(fr1.Header.slice)
			if lenH != len(fr2.Header.slice) {
				log.Errorf("CompareFrames: lengh is not similar: %d/%d", lenH, len(fr2.Header.slice))
				return false
			}

			for i := 0; i < lenH-1; i++ {
				if fr2.Header.slice[i] != fr1.Header.slice[i] {
					log.Errorf("CompareFrames: %s!=%s\n", fr1.Header.slice[i], fr2.Header.slice[i])
					return false
				}
			}
			return true
		}
	}
	return false
}
