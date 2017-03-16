package transport

import (
	"fmt"
	"time"

	"github.com/ventu-io/slf"
)

type LinkDesc struct {
	linkID         string
	address        string
	mode           string
	handler        string
	bufSize        int
	frameProcessor string
}

type Node struct {
	NodeID         string
	LinkDescs      map[string]*LinkDesc
	LinkActives    map[string]*LinkActive
	LinkControls   map[string]*LinkControl
	CommandCh      chan Command
	hasLinks       int
	hasActiveLinks int
	cmdProcessers  []CommandProcesser
	//	Subscribtions  []*chan Message
}

func NewNode() (n *Node) {
	n = &Node{
		LinkDescs:     make(map[string]*LinkDesc),
		LinkControls:  make(map[string]*LinkControl),
		LinkActives:   make(map[string]*LinkActive),
		CommandCh:     make(chan Command),
		cmdProcessers: make([]CommandProcesser, 0),
	}

	dProcesser := &DefaultProcesser{node: n}
	n.cmdProcessers = append(n.cmdProcessers, dProcesser)

	return
}

func (n *Node) AddCmdProcessor(processer CommandProcesser) {
	n.cmdProcessers = append(n.cmdProcessers, processer)
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
			address:        l.Address,
			linkID:         l.LinkID,
			mode:           l.Mode,
			handler:        l.Handler,
			frameProcessor: l.FrameProcessor,
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
	n.CommandCh <- &NodeCommandControlLink{
		NodeCommand: NodeCommand{Cmd: registerControl},
		ctrl:        lControl,
	}
}

func (n *Node) RegisterLinkActive(lActive *LinkActive) {
	log.Debugf("func RegisterLinkActive() %s", lActive.Id())
	n.CommandCh <- &NodeCommandActiveLink{
		NodeCommand: NodeCommand{Cmd: registerActive},
		active:      lActive,
	}
}

func (n *Node) UnregisterLinkControl(lControl *LinkControl) {

	log.Debugf("func UnregisterLinkControl() %s", lControl.getId())
	n.CommandCh <- &NodeCommandControlLink{
		NodeCommand: NodeCommand{Cmd: unregisterControl},
		ctrl:        lControl,
	}
}

func (n *Node) UnregisterLinkActive(lActive *LinkActive) {

	log.Debugf("func UnregisterLinkActive() %s", lActive.Id())
	n.CommandCh <- &NodeCommandActiveLink{
		NodeCommand: NodeCommand{Cmd: unregisterActive},
		active:      lActive,
	}
}

func (n *Node) SendMessage(activeLinkId string, msg string) {
	n.CommandCh <- &NodeCommandSendMessage{
		NodeCommand: NodeCommand{Cmd: sendMessageNode},
		msg:         msg,
	}
}

func closeHelper(closer LinkCloser) {
	closer.Close()
}

func (n *Node) MainLoop() {

	log.Debug("Node.MainLoop() enter")

	var correctCmd bool
	var isExiting bool
	var known bool

	for {
		cmd := <-n.CommandCh
		log.Debugf("MainLoop: get command %+v", cmd)
		correctCmd = false
		for _, processer := range n.cmdProcessers {
			known, isExiting = processer.ProcessCommand(cmd)
			if known {
				correctCmd = true
				break
			}
			if isExiting {
				break
			}
		}
		if !correctCmd {
			log.Errorf("Got unknown command: %v", cmd)
		}
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
	n.CommandCh <- &NodeCommand{
		Cmd: stopNode,
	}

	//TODO: wait with WaitGroup()
	log.Debugf("Node.Stop: waiting active:%d control:%d", n.hasActiveLinks, n.hasLinks)
	for n.hasLinks != 0 || n.hasActiveLinks != 0 {
		time.Sleep(time.Second * time.Duration(1))
		log.Warnf("Node.Stop: waiting active:%d control:%d", n.hasActiveLinks, n.hasLinks)
	}
}
