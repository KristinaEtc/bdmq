package transport

import (
	"fmt"
	"time"

	"github.com/ventu-io/slf"
)

// LinkDesc stores all configs of Links for their launching
type LinkDesc struct {
	linkID         string
	address        string
	mode           string
	handler        string
	bufSize        int
	frameProcessor string
	topic          string
}

// Node is the main entity. It is a datacore of a program
type Node struct {
	NodeID         string                  // ID of Node
	LinkDescs      map[string]*LinkDesc    // Map of all LinkDesc by ID
	LinkActives    map[string]*LinkActive  // Map of all LinkActive by ID
	LinkControls   map[string]*LinkControl // Map of all LinkControl by ID
	CommandCh      chan Command            // A channel for commands
	hasLinks       int
	hasActiveLinks int
	cmdProcessors  []CommandProcessor
	Topics         map[string]map[string]*LinkActive
	//	Subscribtions  map[string]*chan Message	Map of all subscriptions by ID
}

// NewNode creates new instance of Node with DefaultProcesser and CommandProcesser
// and returns it
func NewNode() (n *Node) {
	n = &Node{
		LinkDescs:     make(map[string]*LinkDesc),
		LinkControls:  make(map[string]*LinkControl),
		LinkActives:   make(map[string]*LinkActive),
		CommandCh:     make(chan Command),
		cmdProcessors: make([]CommandProcessor, 0),
		Topics:        make(map[string](map[string]*LinkActive)),
	}

	dProcessor := &DefaultProcessor{node: n}
	n.cmdProcessors = append(n.cmdProcessors, dProcessor)

	return
}

// AddCmdProcessor adds to CommandProcessors a new one.
// Then, when command comes, Node will try to process this in DefaultProcesser;
// if this command couldn't process, it will be sended by turns to all other processer,
// which was added by this command.
func (n *Node) AddCmdProcessor(processor CommandProcessor) {
	n.cmdProcessors = append(n.cmdProcessors, processor)
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

// InitLinkDesc creates LinkDesc for every config struct lDescJSON
func (n *Node) InitLinkDesc(lDescJSON []LinkDescFromJSON) error {

	for _, l := range lDescJSON {

		lDesc := &LinkDesc{
			address:        l.Address,
			linkID:         l.LinkID,
			mode:           l.Mode,
			handler:        l.Handler,
			frameProcessor: l.FrameProcessor,
			topic:          l.Topic,
		}
		n.LinkDescs[l.LinkID] = lDesc
	}

	return nil
}

// InitLinkControl initializes LinkControl and sends a command to Node about new one
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

// RegisterLinkControl sends a command no Node to register a new LinkControl entiry lActive
func (n *Node) RegisterLinkControl(lControl *LinkControl) {

	log.Debugf("func RegisterLinkControl() %s", lControl.getID())
	n.CommandCh <- &NodeCommandControlLink{
		NodeCommand: NodeCommand{Cmd: registerControl},
		ctrl:        lControl,
	}
}

// RegisterLinkActive sends a command no Node to register a new LinkActive entiry lActive
func (n *Node) RegisterLinkActive(lActive *LinkActive) {
	log.Debugf("func RegisterLinkActive() %s", lActive.ID())
	n.CommandCh <- &NodeCommandActiveLink{
		NodeCommand: NodeCommand{Cmd: registerActive},
		active:      lActive,
	}
}

// UnregisterLinkControl sends a command no Node to unregister lControl
func (n *Node) UnregisterLinkControl(lControl *LinkControl) {

	log.Debugf("func UnregisterLinkControl() %s", lControl.getID())
	n.CommandCh <- &NodeCommandControlLink{
		NodeCommand: NodeCommand{Cmd: unregisterControl},
		ctrl:        lControl,
	}
}

// UnregisterLinkActive sends a command no Node to register a new LinkActive entiry
func (n *Node) UnregisterLinkActive(lActive *LinkActive) {

	log.Debugf("func UnregisterLinkActive() %s", lActive.ID())
	n.CommandCh <- &NodeCommandActiveLink{
		NodeCommand: NodeCommand{Cmd: unregisterActive},
		active:      lActive,
	}
}

// RegisterTopic sends a command no Node to register a topic for subscription
func (n *Node) RegisterTopic(topic string, lA *LinkActive) {

	log.Debugf("func  RegisterTopic() %s", topic)
	n.CommandCh <- &NodeCommandTopic{
		NodeCommand: NodeCommand{Cmd: registerTopic},
		topicName:   topic,
		active:      lA,
	}
}

// UnregisterTopic sends a command no Node to unregister a topic for subscription
func (n *Node) UnregisterTopic(topic string, lA *LinkActive) {

	log.Debugf("func  UnregisterTopic() %s", topic)
	n.CommandCh <- &NodeCommandTopic{
		NodeCommand: NodeCommand{Cmd: unregisterTopic},
		topicName:   topic,
		active:      lA,
	}
}

// SendMessage sends a command no Node to send message msg to ActiveLink with ID activeLinkID
func (n *Node) SendMessage(activeLinkID string, msg string) {
	n.CommandCh <- &NodeCommandSendMessage{
		NodeCommand: NodeCommand{Cmd: sendMessageNode},
		msg:         msg,
	}
}

func closeHelper(closer linkCloser) {
	closer.Close()
}

// MainLoop contain a loop which takes Node commands and sends it to Processors
func (n *Node) MainLoop() {

	log.Debug("Node.MainLoop() enter")

	var correctCmd bool
	var isExiting bool
	var known bool

	for {
		cmd := <-n.CommandCh
		correctCmd = false
		for _, processer := range n.cmdProcessors {
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

// Run launches MainLoop() and Initializes Initializes LinkControl in InitLinkControl()
func (n *Node) Run() error {

	log.Debug("Node.Run() enter")

	if n.LinkDescs == nil && len(n.LinkDescs) == 0 {
		log.Debug(errEmptyLinkRepository.Error())
		return errEmptyLinkRepository
	}

	go n.MainLoop()

	for _, lD := range n.LinkDescs {
		go n.InitLinkControl(lD)
	}

	log.Debug("Node.Run() exit")

	return nil
}

// Stop sends stop command for all Links and waites their completion
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
