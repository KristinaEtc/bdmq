package transport

import "github.com/ventu-io/slf"

//HandlerFunc is not implemented
type HandlerFunc string

//Handlers is not implemented
type Handlers map[string]HandlerFunc

var log = slf.WithContext("transport.go")

type creater interface {
	InitNodes([]byte) ([]Noder, error)
	InitHandlers(Handlers) error
}

// Factory is internal class of class, which implements
// creater interface. It contain Init method, which used for
// creating a new node, adding Handlers and opening new transport connection.
type Factory struct {
}

// Init is a method, which used for
// creating a new node, adding Handlers and open new transport connection.
func (f *Factory) Init(c creater, confJSON []byte, h Handlers) ([]Noder, error) {

	nodes, err := c.InitNodes(confJSON)
	if err != nil {
		return nil, err
	}
	c.InitHandlers(h)
	for _, n := range nodes {
		n.RunConn()
	}

	return nodes, nil
}

// Noder is an interface which present a struct with
// common options and list of active connections (linbks;
// implemented by Linker interface).
type Noder interface {
	RunConn()
}

// Linker is an interface which present connection node and
// it's options.
type Linker interface {
	Write(b []byte) (int, error)
	//Disconnect()
	//GetStatus() string
}
