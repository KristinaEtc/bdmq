package transport

import "github.com/ventu-io/slf"

//HandlerFunc is not implemented
type HandlerFunc string

//Handlers is not implemented
type Handlers map[string]HandlerFunc

var log = slf.WithContext("transport.go")

// Noder interface is a struct which implements network endpoint.
type Noder interface {
	InitLinkDesc([]byte) error
	AddLinkDesc([]byte) error
	SendMessage(string, string) error
	Run() error
}

// Linker is an interface which present connection node and
// it's options.
type Linker interface {
	Write(b []byte) (int, error)
	//Disconnect()
	//GetStatus() string
}

type Handler interface {
	OnWrite()
	//...
}
