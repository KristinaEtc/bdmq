package transport

import (
	"net"
	"strings"
	"time"
)

const network = "tcp"

const backOffLimit = time.Duration(time.Second * 600)

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
	bufSize int
}

// Node is a struct which combine connections with common Node's options.
type Node struct {
	NodeID      string
	LinkDesc    map[string]*LinkDesc
	LinkActives map[string]LinkerActive
	LinkStubs   map[string]LinkerActive
}

// NewNode creates an instance of Node struct.
func NewNode() (n *Node) {
	n = &Node{
		LinkDesc:    make(map[string]*LinkDesc),
		LinkActives: make(map[string]LinkerActive),
		LinkStubs:   make(map[string]LinkerActive),
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

// Run reads links and creates nodes with their options.
func (n *Node) Run() error {

	log.Debug("func Run")

	//var wg sync.WaitGroup

	//	var err error

	if n.LinkDesc == nil && len(n.LinkDesc) == 0 {
		log.Debug(ErrEmptyLinkRepository.Error())
		return ErrEmptyLinkRepository
	}

	for _, lD := range n.LinkDesc {

		switch strings.ToLower(lD.mode) {
		case "client":
			//	wg.Add(1)
			go n.initClientLink(lD)
		case "server":
			//	wg.Add(1)
			go n.initServerLink(lD)
		default:
			log.Error(strings.ToLower(lD.mode))
			log.Warnf("%s (ID=%s)\n", ErrWrongNodeMode.Error(), lD.linkID)
		}
	}

	//	wg.Wait()
	log.Debug("func Run closing")
	return nil
}

func (n *Node) initServerLink(linkD *LinkDesc) {

	log.Debug("InitServerLink")

	var err error
	var ln net.Listener

	var secToRecon = time.Duration(time.Second * 2)
	var numOfRecon = 0

	ch := make(chan string)

	n.InitLinkActiveStub(linkD, &ch)

	for {
		ln, err = net.Listen(network, linkD.address)
		if err == nil {
			log.WithField("ID=", linkD.linkID).Debugf("Created a connection with: %s", ln.Addr().String())
			break
		}

		log.WithField("ID=", linkD.linkID).Errorf("Error listen: %s. Reconnecting after %d milliseconds.", err.Error(), secToRecon/1000000.0)
		ticker := time.NewTicker(secToRecon)

		select {
		case _ = <-ticker.C:
			{
				if secToRecon < backOffLimit {

					//randomAdd := int(r1.Float64()*(float64(secToRecon)*0.1) + 0.2*float64(secToRecon))
					randomAdd := secToRecon / 100 * (20 + time.Duration(r1.Int31n(10)))
					//log.Debugf("Random addition=%d", randomAdd/1000000)
					secToRecon = secToRecon*2 + time.Duration(randomAdd)
					//log.Debugf("secToRecon=%d", secToRecon/1000000)
					numOfRecon++
				}
				ticker = time.NewTicker(secToRecon)
				continue
			}

		case command := <-*(n.LinkStubs[linkD.linkID]).(*LinkActiveStub).commandCh:
			{
				if strings.ToLower(command) == commandQuit {
					log.WithField("ID=", linkD.linkID).Info("Got quit command. Closing link.")
					n.LinkStubs[linkD.linkID].Disconnect()
					return
				}
				log.WithField("ID=", linkD.linkID).Warnf("Got impermissible command %s. Ignored.", command)
			}
		}
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

		// race condition
		n.InitLinkActiveWork(linkD, &conn, ((n.LinkStubs[linkD.linkID]).(*LinkActiveStub)).commandCh)
		delete(n.LinkStubs, linkD.linkID)

	}

}

func (n *Node) initClientLink(linkD *LinkDesc) {

	log.Debug("InitClientLink")
}

func (n *Node) InitLinkActiveWork(linkD *LinkDesc, conn *net.Conn, commandCh *chan string) {

	log.Debug("InitLinkActive")

	newActiveLink := LinkActiveWork{
		conn:         conn,
		LinkConf:     linkD,
		LinkActiveID: linkD.linkID + ":" + (*conn).RemoteAddr().String(),
		commandCh:    commandCh,
	}

	h := handlers[linkD.handler].InitHandler(&newActiveLink, n)
	newActiveLink.handler = &h
	n.LinkActives[newActiveLink.LinkActiveID] = &newActiveLink

	go func() {
		(*newActiveLink.handler).OnConnect()
		newActiveLink.Read()
	}()

	//wg.WaitGroup(1)
	log.Debug("InitLinkActive closing")
}

func (n *Node) InitLinkActiveStub(linkD *LinkDesc, commandCh *chan string) {

	log.Debug("InitLinkActiveStub")

	newActiveLink := LinkActiveStub{
		commandCh:    commandCh,
		LinkActiveID: linkD.linkID + ":" + linkD.address,
	}

	n.LinkStubs[linkD.linkID] = &newActiveLink
	log.Debug("InitLinkActiveStub closing")

}

func (n *Node) Stop() {
	log.Debug("Stop(). Not implemented.")
}
