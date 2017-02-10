package transport

import (
	"math/rand"
	"time"
)

var network = "tcp"

const backOffLimit = time.Duration(time.Second * 600)

//for generating reconnect endpoint
var s1 = rand.NewSource(time.Now().UnixNano())
var r1 = rand.New(s1)

type LinkCloser interface {
	Close()
}

type LinkWriter interface {
	Write([]byte) error
	Close()
	Id() string
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
