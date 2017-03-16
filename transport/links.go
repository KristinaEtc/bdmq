package transport

import (
	"math/rand"
	"net"
	"time"
)

var network = "tcp"

const backOffLimit = time.Duration(time.Second * 600)

//for generating reconnect time
var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

// LinkWriter is an interface for LinkActive. It used in handler
// in order to eliminate the dependency on the type LinkActive in other packages
// where Handlers initialized.
type LinkWriter interface {
	Mode() int          // returns Mode of a link; must be server of client
	Write([]byte) error // method for writing to Link
	Close()             // close Link
	ID() string         // returns Link ID
	Conn() net.Conn     // returns conn of Link
}

// linkCloser is an interface which used for gr0ss way to implement graceful closing of ActiveLinks
type linkCloser interface {
	Close()
}

/*
type LinkControl interface {
	InitLinkActive(net.Conn)
	Close()
	NotifyError(error)
	getId() string
	getLinkDesc() *LinkDesc
	getChannel() chan cmdContrlLink
	getNode() *Node
	//	NotifyErrorFromActive(error)
}*/
