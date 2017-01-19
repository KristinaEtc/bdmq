package transport

import "github.com/ventu-io/slf"

var pwdCurr = "KristinaEtc/transport.go"
var log = slf.WithContext(pwdCurr)

/*
//------
// Handler is an interface for communication WithContext
// upper libraries
type Handler interface {
	OnConnect()
	OnDisconnect()
	OnWriter()
	OnError()
}

// File: handler.go
// Handler's realization prototype
//
type Stomp struct {
	maxFramesInPackage int
}

func (s *Stomp) OnConnect() {

}

func (s *Stomp) OnDisconnect() {

}

func (s *Stomp) OnError() {

}

func (s *Stomp) OnWriter() {

}

*/
