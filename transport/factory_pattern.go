package transport

//Handler is not implemented
type Handler string

type creater interface {
	CreateNode() Noder
	InitHandler(Handler)
	RunConn(Noder)
}

// Factory is internal class of class, which implements
// creater interface. It contain Init method, which used for
// creating a new node and opening new transport connection.
type Factory struct {
}

// Init is a method, which used for
// creating a new node and open new transport connection.
func (f *Factory) Init(c creater, nodeConf []byte, h Handler) Noder {
	node := c.CreateNode()
	c.InitHandler(h)
	c.RunConn(node)
	return node
}

// Noder is an interface which implements network transport layer
// and it's methods.
type Noder interface {
	Write(b []byte) (int, error)
	//Disconnect()
	//GetStatus() string
}
