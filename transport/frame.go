package transport

/*
// Frame is an interface for different types of frames
type Frame interface {
	Dump() string
	F()
}
*/

/*
type frameProcessorFactories map[string]FrameProcessorFactory

var frameProcessors frameProcessorFactories = make(map[string]FrameProcessorFactory)

// FrameProcessor is an interface for FrameProcessors.
type FrameProcessor interface {
	Read() error               // read method
	ToByte(interface{}) []byte // from Frame (depends on FrameProcessor) to []byte
}

// FrameProcessorFactory is an interface for creating new FrameProcessor.
type FrameProcessorFactory interface {
	InitFrameProcessor(LinkActive, io.Reader, io.Writer) FrameProcessor // creates new FrameProcessor
}

// RegisterFrameProcessorFactory added FrameProcessorFactory fFactory with name frameProcesserName.
// Further this FrameProcessor could be used with this LinkActives.
func RegisterFrameProcessorFactory(frameProcesserName string, fFactory FrameProcessorFactory) {
	log.Debug("func RegisterFrameFactory()")
	frameProcessors[frameProcesserName] = fFactory
}

//-------------------------------------------------------
//	 defaultFrame
//-------------------------------------------------------

// defaultFrameProcessorFactory is a factory for creating defaultFrameProcessor.
type defaultFrameProcessorFactory struct {
}

var dFrameProcessorFactory defaultFrameProcessorFactory

// defaultFrameProcessor implements FrameProcesser interface.
// Used by default.
type defaultFrameProcessor struct {
	linkActive LinkActive
	reader     bufio.Reader
	writer     io.Writer
	log        slf.Logger
	handler    Handler
}

// InitFrameProcessor creates a new entity of defaultFrameProcessor and adds it to Node process slice.
func (d defaultFrameProcessorFactory) initFrameProcessor(lActive LinkActive, rd net.Conn, wd net.Conn) FrameProcessor {

	dFrameProcesser := &defaultFrameProcessor{
		writer:  wd,
		log:     slf.WithContext("defaultFrameProcessor").WithFields(slf.Fields{"ID": lActive.ID()}),
		handler: lActive.getHandler(),
	}

	dFrameProcesser.newReader(rd)

	return dFrameProcesser
}

func (d *defaultFrameProcessor) newReader(rd io.Reader) (r io.Reader) {
	r = bufio.NewReader(rd)
	return
}

/*
// Read waiting input data from reader and then call
// handler.OnRead() to process it.
func (d *defaultFrameProcessor) Read() error {
	for {
		message, err := d.reader.ReadBytes('\n')
		if err != nil {
			d.log.Errorf("Error read: %s", err.Error())
			return err
		}
		d.handler.OnRead(message)
	}
}

// ToByte converts Frame (for different FrameProcessers different types) type to slices of bytes.
func (d *defaultFrameProcessor) ToByte(msg interface{}) []byte {
	return (msg.([]byte))
}

*/
