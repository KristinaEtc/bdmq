package stomp

/*
import (
	"io"

	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

type StompFrameProcessorFactory struct {
}

var FrameProcessorFactory StompFrameProcessorFactory

// StompFrameProcessor implements FrameProcesser interface.
type StompFrameProcessor struct {
	linkActiveID string
	reader       *frame.Reader
	writer       *frame.Writer
	log          slf.Logger
}

func (h StompFrameProcessorFactory) InitFrameProcessor(rd io.Reader, wd io.Writer, log slf.Logger) transport.FrameProcessor {

	log.Debugf("InitDefaultFrame")
	sFrameProcesser := &StompFrameProcessor{
		writer: frame.NewWriter(wd),
		reader: frame.NewReader(rd),
		log:    log,
	}

	return sFrameProcesser
}

func (s *StompFrameProcessor) Read() (message []byte, err error) {

	for {
		ff, err := s.reader.Read()
		if err != nil {
			if err == io.EOF {
				s.log.Errorf("connection closed: eof")
			} else {
				s.log.Errorf("read failed: %s", err.Error())
			}
			break
		}
		if ff == nil {
			log.Infof("heartbeat")
			continue
		}
		return ff.Body, err
	}
	return nil, err
}

func (s *StompFrameProcessor) Write(msg []byte) {
	var frame = frame.New(
		"SEND", "test1", "2", "message", string(msg))

	s.writer.Write(frame)
}
*/
