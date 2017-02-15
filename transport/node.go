package transport

import (
	"fmt"
	"time"

	"github.com/ventu-io/slf"
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
	LinkControls   map[string]*LinkControl
	commandCh      chan *NodeCommand
	hasLinks       int
	hasActiveLinks int
}

func NewNode() (n *Node) {
	n = &Node{
		LinkDescs:    make(map[string]*LinkDesc),
		LinkControls: make(map[string]*LinkControl),
		LinkActives:  make(map[string]*LinkActive),
		commandCh:    make(chan *NodeCommand),
	}
	return
}

func checkLinkMode(mode string) (int, error) {
	switch mode {
	case "client":
		return 0, nil
	case "server":
		return 1, nil
	}
	return 2, fmt.Errorf("Wrong link mode: %s", mode)
}

func (n *Node) InitLinkDesc(lDescJSON []LinkDescFromJSON) error {

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

func (n *Node) InitLinkControl(lD *LinkDesc) {

	log := slf.WithContext("LinkControl").WithFields(slf.Fields{"ID": lD.linkID})
	log.Debugf("func InitLinkControl() %+v", lD)

	mode, err := checkLinkMode(lD.mode)
	if err != nil {
		log.Errorf("checkLinkMode error %s", err.Error())
		return
	}

	linkControl := &LinkControl{
		linkDesc:  lD,
		node:      n,
		commandCh: make(chan cmdContrlLink),
		log:       log,
		mode:      mode,
	}

	n.RegisterLinkControl(linkControl)
	defer n.UnregisterLinkControl(linkControl)

	switch mode {
	case 0:
		linkControl.WorkClient()
	case 1:
		linkControl.WorkServer()
	}

	log.Debugf("func InitLinkControl() %+v closing", lD)
}

func (n *Node) RegisterLinkControl(lControl *LinkControl) {

	log.Debugf("func RegisterLinkControl() %s", lControl.getId())
	n.commandCh <- &NodeCommand{
		cmd:  registerControl,
		ctrl: lControl,
	}
}

func (n *Node) RegisterLinkActive(lActive *LinkActive) {
	log.Debugf("func RegisterLinkActive() %s", lActive.Id())
	n.commandCh <- &NodeCommand{
		cmd:    registerActive,
		active: lActive,
	}
}

func (n *Node) UnregisterLinkControl(lControl *LinkControl) {

	log.Debugf("func UnregisterLinkControl() %s", lControl.getId())
	n.commandCh <- &NodeCommand{
		cmd:  unregisterControl,
		ctrl: lControl,
	}
}

func (n *Node) UnregisterLinkActive(lActive *LinkActive) {

	log.Debugf("func UnregisterLinkActive() %s", lActive.Id())
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

func closeHelper(closer LinkCloser) {
	closer.Close()
}

func (n *Node) processCommand(cmdMsg *NodeCommand) (isExiting bool) {
	log.Debugf("process command=%s", cmdMsg.cmd.String())
	switch cmdMsg.cmd {
	case registerActive:
		{
			n.LinkActives[cmdMsg.active.Id()] = cmdMsg.active
			//n.hasActiveLinks = len(n.LinkActives) > 0
			n.hasActiveLinks++
			log.Debugf("[registerActive] linkA=%d, links=%d", n.hasActiveLinks, n.hasLinks)
			return false
		}
	case unregisterActive:
		{
			if len(n.LinkActives) > 0 {
				for _, lA := range n.LinkActives {
					log.Debugf("%s", lA.Id())
					//return false
				}
			} else {
				log.Debug("NIL")
			}

			delete(n.LinkActives, cmdMsg.active.Id())
			if n.hasActiveLinks != 0 {
				n.hasActiveLinks--
			} else {
				log.Debug("h.hasActiveLinks<0")
			}
			log.Debugf("[unregisterActive] linkA=%d, links=%d", n.hasActiveLinks, n.hasLinks)

			return n.hasActiveLinks == 0 && n.hasLinks == 0

			//n.hasActiveLinks = len(n.LinkActives) > 0
			//return !(n.hasLinks || n.hasActiveLinks)
			//return false
		}
	case registerControl:
		{
			n.LinkControls[cmdMsg.ctrl.getId()] = cmdMsg.ctrl
			//n.hasLinks = len(n.LinkControls) > 0
			n.hasLinks++
			log.Debugf("[registerControl] linkA=%d, links=%d", n.hasActiveLinks, n.hasLinks)
			return false
		}
	case unregisterControl:
		{
			delete(n.LinkControls, cmdMsg.ctrl.getId())
			//	n.hasLinks = len(n.LinkControls) > 0
			//	return !(n.hasLinks || n.hasActiveLinks)
			if n.hasLinks != 0 {
				n.hasLinks--
			} else {
				log.Debug("h.hasLinks = 0!!")
			}
			log.Debugf("[unregisterControl] linkA=%d, links=%d", n.hasActiveLinks, n.hasLinks)
			return n.hasActiveLinks == 0 && n.hasLinks == 0
		}
	case stopNode:
		{
			log.Infof("Stop received")
			if len(n.LinkActives) > 0 {
				for linkID, lA := range n.LinkActives {
					log.Debugf("Send close to active link %s %s", linkID, lA.Id())
					go closeHelper(lA)
				}
			}
			if len(n.LinkControls) > 0 {
				for linkID, lC := range n.LinkControls {
					log.Debugf("Send close to link control %s %s", linkID, lC.getId())
					go closeHelper(lC)
				}
			}
			return false
		}
	case sendMessageNode:
		{
			if len(n.LinkActives) > 0 {
				for _, lA := range n.LinkActives {
					log.Debug("for testing i'm choosing all active links")
					lA.SendMessage(cmdMsg.msg)
					//return false
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

	log.Debug("Node.MainLoop() enter")

	for {
		cmd := <-n.commandCh
		isExiting := n.processCommand(cmd)
		if isExiting {
			break
		}
	}
	log.Debug("Node.MainLoop() exit")
}

func (n *Node) Run() error {

	log.Debug("Node.Run() enter")

	if n.LinkDescs == nil && len(n.LinkDescs) == 0 {
		log.Debug(ErrEmptyLinkRepository.Error())
		return ErrEmptyLinkRepository
	}

	go n.MainLoop()

	for _, lD := range n.LinkDescs {
		go n.InitLinkControl(lD)
	}

	log.Debug("Node.Run() exit")

	return nil
}

func (n *Node) Stop() {
	log.Debug("Node.Stop()")

	//todo: add checking if channel is exist
	n.commandCh <- &NodeCommand{
		cmd: stopNode,
	}

	//TODO: wait with WaitGroup()
	log.Debugf("Node.Stop: waiting active:%d control:%d", n.hasActiveLinks, n.hasLinks)
	for n.hasLinks != 0 || n.hasActiveLinks != 0 {
		time.Sleep(time.Second * time.Duration(1))
		log.Warnf("Node.Stop: waiting active:%d control:%d", n.hasActiveLinks, n.hasLinks)
	}
}
