package transport

import (
	"math/rand"
	"net"
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
	InitLinkActive(net.Conn)
	Close()
	NotifyError(error)
	getId() string
	getLinkDesc() *LinkDesc
	getChannel() chan cmdContrlLink
	getNode() *Node
	//	NotifyErrorFromActive(error)
}

func initLinkActive(lCntl LinkControl, conn net.Conn) {

	linkActive := LinkActive{
		conn:         conn,
		linkDesc:     lCntl.getLinkDesc(),
		LinkActiveID: lCntl.getId() + ":" + conn.RemoteAddr().String(),
		commandCh:    make(chan cmdActiveLink),
		linkControl:  lCntl,
	}

	log.Debug(lCntl.getLinkDesc().handler)

	ch := lCntl.getChannel()

	if _, ok := handlers[lCntl.getLinkDesc().handler]; !ok {
		log.Error("No handler! Closing linkControl")
		ch <- cmdContrlLink{
			cmd: errorControlLink,
			err: "Error: " + "No handler! Closing linkControl",
		}
	}

	node := lCntl.getNode()
	h := handlers[lCntl.getLinkDesc().handler].InitHandler(&linkActive, node)
	linkActive.handler = h

	node.RegisterLinkActive(&linkActive)

	go linkActive.WaitCommand(conn)
	linkActive.handler.OnConnect()
	linkActive.Read()
	log.Debug("initLinkActive exiting")
	//linkActive.handler.OnDisconnect()

	node.UnregisterLinkActive(&linkActive)
}
