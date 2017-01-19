package transport

//Handler is not implemented
type Handler string

type creater interface {
	CreateNode() Noder
	InitHandler(Handler)
	RunConn(Noder)
}

type Factory struct {
}

func (f *Factory) Init(c creater, nodeConf []byte, h Handler) Noder {
	node := c.CreateNode()
	c.InitHandler(h)
	c.RunConn(node)
	return node
}

type Noder interface {
	Disconnect()
	Write(string)
	GetStatus() string
}
