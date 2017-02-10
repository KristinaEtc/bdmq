package transport

import "time"
import "fmt"

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
	hasLinks       int
	hasActiveLinks int
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

func (n *Node) InitLinkControl(lD *LinkDesc) {
	log.Debug("func InitLinkControl()")

	linkControl := LinkControl{
		linkDesc:  lD,
		node:      n,
		commandCh: make(chan cmdContrlLink),
	}

	mode, err := checkLinkMode(lD.mode)
	if err != nil {
		log.Error(err.Error())
		return
	}

	n.RegisterLinkControl(linkControl)
	defer n.UnregisterLinkControl(linkControl)

	switch mode {
	case 0:
		linkControl.WorkClient()
	case 1:
		linkControl.WorkServer()
	}

	log.Debug("func InitServerLinkControl() closing")
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
					log.Debugf("Send close to %s", linkID)
					go lA.Close()
				}
			}
			if len(n.LinkControls) > 0 {
				for linkID, lC := range n.LinkControls {
					log.Debugf("Send close to %s", linkID)
					go lC.Close()
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
		go n.InitLinkControl(lD)
	}

	log.Debug("func Run() closing")

	return nil
}

func (n *Node) Stop() {
	log.Debug("func Stop()")

	//todo: add checking if channel is exist
	n.commandCh <- &NodeCommand{
		cmd: stopNode,
	}

	//TODO: wait with WaitGroup()
	log.Warn("Waiting")
	log.Debugf("hasLinks=%d, hasActive=%d", n.hasLinks, n.hasActiveLinks)
	for n.hasLinks != 0 || n.hasActiveLinks != 0 {
		if n.hasActiveLinks != 0 {
			log.Debugf("hasActive=%d", n.hasActiveLinks)
		}
		if n.hasLinks != 0 {
			log.Debugf("hasLinks=%d", n.hasLinks)
		}
		time.Sleep(time.Second * time.Duration(1))
		log.Warnf("waiting")
	}
}
