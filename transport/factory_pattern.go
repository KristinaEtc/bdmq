package transport

import "github.com/ventu-io/slf"
import "fmt"

//HandlerFunc & Handler are not implemented
type HandlerFunc string
type Handler map[string]HandlerFunc

var log = slf.WithContext("transport.go")

type factoryNodeCreater interface {
	ParseNodeConfig([]byte) (Noder, error)
	AddHandlers(Noder, Handler) error
	RunConn(Noder) error
}

// Factory is internal class of class, which implements
// creater interface. It contain Init method, which used for
// creating a new node, adding Handlers and opening new transport connection.
type Factory struct {
}

// Init is a method, which used for
// creating a new node, adding Handlers and open new transport connection.
func (f *Factory) Init(c creater, nodeConf []byte, h Handler) (Noder, error) {

	nodes, err := c.ParseNodeConfig(nodeConf)
	if err != nil {
		err = fmt.Errorf("Error while parsing nodes' config: %s", err.Error())
		return nil, err
	}

	err = c.AddHandlers(nodes, h)
	if err != nil {
		err = fmt.Errorf("Error while adding handlers: %s", err.Error())
		return nil, err
	}

	err = c.RunConn(nodes)
	if err != nil {
		err = fmt.Errorf("Error while start connections: %s", err.Error())
		return nil, err
	}

	return nodes, nil
}

// Noder is an interface which implements network transport layer
// and it's methods.
type Noder interface {
	Write(b []byte) (int, error)
	//Disconnect()
	//GetStatus() string
}
