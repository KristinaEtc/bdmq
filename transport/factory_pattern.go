package transport

//Handler is not implemented
type Handler string

type creater interface {
	CreateNode() noder
	InitHandler(Handler)
	RunConn(string)
}

type Factory struct {
}

func (f *Factory) Init(c creater, typeOfNode string) noder {
	node := c.CreateNode()

	//node.RunConn(typeOfNode)
	return node
}

type noder interface {
	Disconnect()
	Write()
	GetStatus()
}
