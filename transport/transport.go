package transport

import (
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
	NodeID   string
	LinkDesc map[string]*LinkDesc
	//LinkActives  map[string]LinkActiver
	//LinkControls map[string]LinkActiver
	LinkControls map[string]LinkControl
	CommandCh    chan *NodeCommand
	hasLinks     bool
}

type NodeCommand struct {
	command string
	ctrl    LinkControl
}

// NewNode creates an instance of Node struct.
func NewNode() (n *Node) {
	n = &Node{
		LinkDesc: make(map[string]*LinkDesc),
		//LinkActives:  make(map[string]LinkActiver),
		LinkControls: make(map[string]LinkControl),
		CommandCh:    make(chan *NodeCommand),
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

	go n.MainLoop()

	for _, lD := range n.LinkDesc {

		switch strings.ToLower(lD.mode) {
		case "client":
			//	wg.Add(1)
			go n.initClientLinkControl(lD)
		case "server":
			//	wg.Add(1)
			go n.initServerLinkControl(lD)
		default:
			log.Error(strings.ToLower(lD.mode))
			log.Warnf("%s (ID=%s)\n", ErrWrongNodeMode.Error(), lD.linkID)
		}
	}

	//	wg.Wait()
	log.Debug("func Run closing")
	return nil
}

func (n *Node) MainLoop() {
	log.Debug("func MainLoop")
	for {
		cmd := <-n.CommandCh
		isExiting := n.processCommand(cmd)
		if isExiting {
			break
		}
		//todo: wait quit
	}
	log.Debug("MainLoop exiting")
}

func (n *Node) processCommand(cmd *NodeCommand) bool {
	log.Debugf("processCommand %s", cmd.command)
	if cmd.command == "stop" {
		log.Infof("stop received")
		for k, v := range n.LinkControls {
			log.Infof("send close to %s", k)
			v.Close()
		}
		return false
	}
	if cmd.command == "unregister" {
		delete(n.LinkControls, cmd.ctrl.Id())
		n.hasLinks = len(n.LinkControls) > 0
		return !n.hasLinks
	}
	if cmd.command == "register" {
		n.LinkControls[cmd.ctrl.Id()] = cmd.ctrl
		n.hasLinks = len(n.LinkControls) > 0
		return false
	}
	log.Warnf("unknown command %s", cmd.command)
	return false
}

func (n *Node) RegisterLinkControl(linkControl LinkControl) {
	n.CommandCh <- &NodeCommand{
		command: "register",
		ctrl:    linkControl,
	}
}

func (n *Node) UnregisterLinkControl(linkControl LinkControl) {
	n.CommandCh <- &NodeCommand{
		command: "unregister",
		ctrl:    linkControl,
	}
}

func (n *Node) initServerLinkControl(linkD *LinkDesc) {

	log.Debug("initServerLinkControl")

	linkControl := LinkControlServer{
		CommandCh: make(chan string),
		//LinkActiveID: linkD.linkID + ":" + linkD.address,
		Node:     n,
		LinkDesc: linkD,
	}

	//n.LinkControls[linkD.linkID] = &newLinkControlS
	n.RegisterLinkControl(&linkControl)

	for {
		ln, err := linkControl.Listen()
		if err != nil {
			log.Errorf("Listen error %s", err.Error())
			break
		}
		//run accept goroutine
		go linkControl.Accept(ln)
		isExiting := linkControl.WaitCommand(ln)
		if isExiting {
			break
		}
	}

	n.UnregisterLinkControl(&linkControl)

}

func (n *Node) initClientLinkControl(linkD *LinkDesc) {
	log.Debug("InitClientLink")
}

/*
func (n *Node) InitLinkActive(linkD *LinkDesc, conn *net.Conn, node *Node) {

	log.Debug("InitLinkActive")

	newActiveLink := LinkActive{
		conn:         conn,
		LinkConf:     linkD,
		LinkActiveID: linkD.linkID + ":" + (*conn).RemoteAddr().String(),
		commandCh:    make(chan string),
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
}*/

func (n *Node) Stop() {
	log.Debugf("Stop")
	n.CommandCh <- &NodeCommand{
		command: "stop",
	}
	//TODO: wait with WaitGroup
	log.Warnf("waiting")
	for n.hasLinks {
		time.Sleep(time.Second * time.Duration(1))
		log.Warnf("waiting")
	}
}
