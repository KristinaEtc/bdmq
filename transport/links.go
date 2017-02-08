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

type LinkWriter interface {
	Write(string) error
	Close()
	Id() string
}

type LinkControl interface {
	Close()
	Id() string
}
