package transport

import (
	"bufio"
	"net"
	"strings"
	"sync"
	"time"
)

const network = "tcp"

const backOffLimit = time.Duration(time.Minute * 10)

// status enum
const (
	statusOnline       = iota
	statusOffline      = iota
	statusReconnecting = iota
)

// LinkDesc contains options for connecting to an address.
// It is store options for creating connection for ctrtain connection.
type LinkDesc struct {
	linkID  string
	address string
	mode    string //client/server
	handler string
}

// LinkActive provides a TCP connection with auto-reconnect capabilities.
type LinkActive struct {
	conn         *net.Conn
	lock         sync.RWMutex
	status       int32
	LinkConf     *LinkDesc
	LinkActiveID string
	handler      *Handler
	//parser       *parser.Parser
}

// Node is a struct which combine connections with common Node's options.
type Node struct {
	NodeID      string
	LinkDesc    map[string]*LinkDesc
	LinkActives map[string]*LinkActive
}

//--------------------------------------

// NewNode creates an instance of Node struct.
func NewNode() (n *Node) {
	n = &Node{
		LinkDesc:    make(map[string]*LinkDesc),
		LinkActives: make(map[string]*LinkActive),
	}
	return
}

// InitLinkDesc parses config with nodes' options and creates nodes with such
// options.
// InitLinkDesc is a method of Noder interface from transport library.
func (n *Node) InitLinkDesc(lDescJSON []LinkDescFromJSON) error {

	log.Debug("InitLinkDesc")

	for _, l := range lDescJSON {

		lDesc := &LinkDesc{
			address: l.Address,
			linkID:  l.LinkID,
			mode:    l.Mode,
			handler: l.Handler,
		}

		n.LinkDesc[l.LinkID] = lDesc
	}

	return nil
}

//Run reads links and creates nodes with their options.
func (n *Node) Run() error {

	log.Debug("func Run")

	var wg sync.WaitGroup

	//	var err error

	if n.LinkDesc == nil && len(n.LinkDesc) == 0 {
		log.Debug(ErrEmptyLinkSlice.Error())
		return ErrEmptyLinkSlice
	}

	for _, lD := range n.LinkDesc {

		switch strings.ToLower(lD.mode) {
		case "client":
			wg.Add(1)
			go n.initClientLink(lD)
		case "server":
			wg.Add(1)
			go n.initServerLink(lD)
		default:
			log.Error(strings.ToLower(lD.mode))
			log.Warnf("%s (ID=%s)\n", ErrWrongNodeMode.Error(), lD.linkID)
		}
	}

	wg.Wait()
	log.Debug("func Run closing")
	return nil
}

func (n *Node) initServerLink(linkD *LinkDesc) {

	log.Debug("InitServerLink")

	var err error
	var ln net.Listener
	// TODO: change this block when refactoring
	// listen on all interfaces
	for {
		ln, err = net.Listen(network, linkD.address)
		if err != nil {
			log.WithField("ID=", linkD.linkID).Errorf("Error listen: %s", err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		break
	}

	for {

		// accept connection on port
		conn, err := ln.Accept()
		if err != nil {
			log.WithField("ID=", linkD.linkID).Errorf("Error accept: %s", err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		log.WithField("link ID=", linkD.linkID).Debug("New client")
		n.InitLinkActive(linkD, &conn)

	}

}

func (n *Node) InitLinkActive(linkD *LinkDesc, conn *net.Conn) {

	log.Debug("InitLinkActive")

	newActiveLink := LinkActive{
		conn:         conn,
		LinkConf:     linkD,
		LinkActiveID: linkD.linkID + ":" + (*conn).RemoteAddr().String(),
	}

	h := handlers[linkD.handler].InitHandler(&newActiveLink, n)
	newActiveLink.handler = &h
	n.LinkActives[newActiveLink.LinkActiveID] = &newActiveLink

	go (*newActiveLink.handler).OnConnect()
	go newActiveLink.Read()

	//	WaitGroup(2)

	log.Debug("InitLinkActive closing")

}

func (n *Node) initClientLink(linkD *LinkDesc) {
	log.Debug("Init client node")

}

func (lA *LinkActive) Read() {
	for {
		message, err := bufio.NewReader(*lA.conn).ReadString('\n')
		if err != nil {
			log.WithField("ID=", lA.LinkActiveID).Errorf("Error read: %s", err.Error())
		}
		log.WithField("ID=", lA.LinkActiveID).Debugf("Message Received: %s", string(message))
		go (*lA.handler).OnRead(message)
	}
}

func (lA *LinkActive) Write(msg string) error {
	(*lA.conn).Write([]byte(msg + "\n"))
	return nil
}
