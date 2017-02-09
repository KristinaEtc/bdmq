package transport

import (
	"strings"
	"time"
)

type LinkDesc struct {
	linkID  string
	address string
	mode    string
	handler string
	bufSize int
}

type Node struct {
	NodeID         string
	LinkDescs      map[string]*LinkDesc
	LinkActives    map[string]*LinkActive
	LinkControls   map[string]LinkControl
	commandCh      chan *NodeCommand
	hasLinks       bool
	hasActiveLinks bool
}

func NewNode() (n *Node) {
	n = &Node{
		LinkDescs:    make(map[string]*LinkDesc),
		LinkControls: make(map[string]LinkControl),
		LinkActives:  make(map[string]*LinkActive),
		commandCh:    make(chan *NodeCommand),
	}
	return
}

func (n *Node) InitLinkDesc(lDescJSON []LinkDescFromJSON) error {

	log.Debug("func InitLinkDesc()")

	for _, l := range lDescJSON {

		lDesc := &LinkDesc{
			address: l.Address,
			linkID:  l.LinkID,
			mode:    l.Mode,
			handler: l.Handler,
		}
		n.LinkDescs[l.LinkID] = lDesc
	}

	return nil
}

func (n *Node) InitServerLinkControl(lD *LinkDesc) {
	log.Debug("func InitServerLinkControl()")

	linkControl := LinkControlServer{
		linkDesc:  lD,
		node:      n,
		commandCh: make(chan cmdContrlLink),
	}

	n.RegisterLinkControl(&linkControl)

	for {
		ln, err := linkControl.Listen()
		if err != nil {
			log.Errorf("Listen error: %s", err.Error())
			break
		}
		go linkControl.Accept(ln)
		isExiting := linkControl.WaitCommand(ln)
		if isExiting {
			break
		}
	}
	n.UnregisterLinkControl(&linkControl)

	log.Debug("func InitServerLinkControl() closing")
}

func (n *Node) InitClientLinkControl(lD *LinkDesc) {
	log.Debug("func InitClientLinkControl()")

	linkControl := LinkControlClient{
		linkDesc:  lD,
		node:      n,
		commandCh: make(chan cmdContrlLink),
	}

	n.RegisterLinkControl(&linkControl)

	for {
		conn, err := linkControl.Dial()
		if err != nil {
			log.Errorf("Dial error: %s", err.Error())
			break
		}
		isExiting := linkControl.WaitCommand(conn)
		if isExiting {
			break
		}
	}
	n.UnregisterLinkControl(&linkControl)

	log.Debug("func InitClientLinkControl() closing")
}

func (n *Node) RegisterLinkControl(lControl LinkControl) {

	log.Debug("func RegisterLinkControl()")
	n.commandCh <- &NodeCommand{
		cmd:  registerControl,
		ctrl: lControl,
	}
}

func (n *Node) RegisterLinkActive(lActive *LinkActive) {
	log.Debug("func RegisterLinkActive()")
	n.commandCh <- &NodeCommand{
		cmd:    registerActive,
		active: lActive,
	}
}

func (n *Node) UnregisterLinkControl(lControl LinkControl) {

	log.Debug("func UnregisterLinkControl()")
	n.commandCh <- &NodeCommand{
		cmd:  unregisterControl,
		ctrl: lControl,
	}
}

func (n *Node) UnregisterLinkActive(lActive *LinkActive) {

	log.Debug("func UnregisterLinkActive()")
	n.commandCh <- &NodeCommand{
		cmd:    unregisterActive,
		active: lActive,
	}
}

func (n *Node) SendMessage(activeLinkId string, msg string) {
	n.commandCh <- &NodeCommand{
		cmd: sendMessageNode,
		msg: msg,
	}
}

func (n *Node) processCommand(cmdMsg *NodeCommand) (isExiting bool) {
	log.Debugf("process command=%s", cmdMsg.cmd.String())
	switch cmdMsg.cmd {
	case registerActive:
		{
			n.LinkActives[cmdMsg.active.Id()] = cmdMsg.active
			n.hasActiveLinks = len(n.LinkActives) > 0
			return false
		}
	case unregisterActive:
		{
			delete(n.LinkActives, cmdMsg.active.Id())
			n.hasActiveLinks = len(n.LinkControls) > 0
			return !(n.hasLinks && n.hasActiveLinks)
		}
	case registerControl:
		{
			n.LinkControls[cmdMsg.ctrl.getId()] = cmdMsg.ctrl
			n.hasLinks = len(n.LinkControls) > 0
			return false
		}
	case unregisterControl:
		{
			delete(n.LinkControls, cmdMsg.ctrl.getId())
			n.hasLinks = len(n.LinkControls) > 0
			return !(n.hasLinks && n.hasActiveLinks)
		}
	case stopNode:
		{
			log.Infof("Stop received")
			if len(n.LinkActives) > 0 {
				for linkID, lA := range n.LinkActives {
					log.Debugf("Send close to %s", linkID)
					lA.Close()
				}
			}
			if len(n.LinkControls) > 0 {
				for linkID, lC := range n.LinkControls {
					log.Debugf("Send close to %s", linkID)
					lC.Close()
				}
			}
			return false
		}
	case sendMessageNode:
		{
			if len(n.LinkActives) > 0 {
				for _, lA := range n.LinkActives {
					log.Debug("for testing i'm choosing the 1th of active links")
					lA.SendMessage(cmdMsg.msg)
					return false
				}
			}
			log.Debug("SendMessage: no active links")
			return false
		}
	default:
		{
			log.Warnf("Unknown command: %s", cmdMsg.cmd)
			return false
		}
	}
}

func (n *Node) MainLoop() {
	log.Debug("func MainLoop()")

	for {
		cmd := <-n.commandCh
		isExitisng := n.processCommand(cmd)
		if isExitisng {
			break
		}
	}
	log.Debug("func MainLoop() closing")
}

func (n *Node) Run() error {

	log.Debug("func Run()")

	if n.LinkDescs == nil && len(n.LinkDescs) == 0 {
		log.Debug(ErrEmptyLinkRepository.Error())
		return ErrEmptyLinkRepository
	}

	go n.MainLoop()

	for _, lD := range n.LinkDescs {
		switch strings.ToLower(lD.mode) {
		case "client":
			go n.InitClientLinkControl(lD)

		case "server":
			go n.InitServerLinkControl(lD)

		default:
			log.Errorf("Wrong link mode: %s; ignored.", lD.mode)
		}
	}

	log.Debug("func Run() closing")

	return nil
}

func (n *Node) Stop() {
	log.Debug("func Stop()")
	n.commandCh <- &NodeCommand{
		cmd: stopNode,
	}

	//TODO: wait with WaitGroup()
	log.Warn("Waiting")
	for n.hasLinks && n.hasActiveLinks {
		if n.hasActiveLinks {
			log.Debug("hasActive")
		}
		if n.hasLinks {
			log.Debug("hasLinks")
		}
		time.Sleep(time.Second * time.Duration(1))
		log.Warnf("waiting")
	}
}
