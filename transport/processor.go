package transport

// Command is an interface of all commands of Nodes
// (for example defaultCommandNode, stompCommands from package stomp)
type Command interface {
	GetCommandID() CommandID
}

// CommandProcessor is an interface which was created for method ProcessCommand.
//
// Processcommand rocesses all Nodes' commands.
type CommandProcessor interface {
	ProcessCommand(Command) (bool, bool)
	CommandToString(CommandID) string
}

// DefaultProcessor is used by Node by default; contains basic commands.
type DefaultProcessor struct {
	node    *Node
	handler Handler
}

// ProcessCommand contains a loop where commands handles
func (d *DefaultProcessor) ProcessCommand(cmd Command) (known bool, isExiting bool) {
	var id = cmd.GetCommandID()
	log.Debugf("Process command [%s]", d.CommandToString(cmd.GetCommandID()))

	n := d.node

	switch id {
	case registerActive:
		{
			cmdActive, ok := cmd.(*NodeCommandActiveLink)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}

			n.LinkActives[cmdActive.active.ID()] = cmdActive.active
			n.hasActiveLinks++
			log.Debugf("[registerActive] linkA=%d, links=%d", n.hasActiveLinks, n.hasLinks)
			return true, false
		}
	case unregisterActive:
		{
			cmdActive, ok := cmd.(*NodeCommandActiveLink)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}

			if len(n.LinkActives) > 0 {
				for _, lA := range n.LinkActives {
					log.Debugf("%s", lA.ID())
					//return false
				}
			} else {
				log.Debug("NIL")
			}

			delete(n.LinkActives, cmdActive.active.ID())
			if n.hasActiveLinks != 0 {
				n.hasActiveLinks--
			} else {
				log.Debug("h.hasActiveLinks<0")
			}
			log.Debugf("[unregisterActive] linkA=%d, links=%d", n.hasActiveLinks, n.hasLinks)

			return true, n.hasActiveLinks == 0 && n.hasLinks == 0

			//n.hasActiveLinks = len(n.LinkActives) > 0
			//return !(n.hasLinks || n.hasActiveLinks)
			//return false
		}
	case registerControl:
		{
			cmdControl, ok := cmd.(*NodeCommandControlLink)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}

			n.LinkControls[cmdControl.ctrl.getID()] = cmdControl.ctrl
			//n.hasLinks = len(n.LinkControls) > 0
			n.hasLinks++
			log.Debugf("[registerControl] linkA=%d, links=%d", n.hasActiveLinks, n.hasLinks)
			return true, false
		}
	case unregisterControl:
		{
			cmdControl, ok := cmd.(*NodeCommandControlLink)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}

			delete(n.LinkControls, cmdControl.ctrl.getID())
			//	n.hasLinks = len(n.LinkControls) > 0
			//	return !(n.hasLinks || n.hasActiveLinks)
			if n.hasLinks != 0 {
				n.hasLinks--
			} else {
				log.Debug("h.hasLinks = 0!!")
			}
			log.Debugf("[unregisterControl] linkA=%d, links=%d", n.hasActiveLinks, n.hasLinks)
			return true, n.hasActiveLinks == 0 && n.hasLinks == 0
		}
	case stopNode:
		{

			log.Infof("Stop received")
			if len(n.LinkActives) > 0 {
				for linkID, lA := range n.LinkActives {
					log.Debugf("Send close to active link %s %s", linkID, lA.ID())
					go closeHelper(lA)
				}
			}
			if len(n.LinkControls) > 0 {
				for linkID, lC := range n.LinkControls {
					log.Debugf("Send close to link control %s %s", linkID, lC.getID())
					go closeHelper(lC)
				}
			}
			return true, false
		}
	case sendStringNode:
		{
			cmdMessage, ok := cmd.(*NodeCommandSendString)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}

			if len(n.LinkActives) > 0 {
				for _, lA := range n.LinkActives {
					log.Debug("for testing i'm choosing all active links")
					lA.SendStringActive([]byte(cmdMessage.msg))
					//return false
				}
			}
			return true, false
		}
	default:
		{
			//log.Warnf("Unknown command: %s", cmdMsg.cmd)
			return false, false
		}
	}
}

// CommandToString returns string representation from commandID c
func (d *DefaultProcessor) CommandToString(c CommandID) string {
	switch c {
	case 0:
		return "registerControl"
	case 1:
		return "unregisterControl"
	case 2:
		return "registerActive"
	case 3:
		return "unregisterActive"
	case 4:
		return "stop"
	case 5:
		return "sendString"
	case 6:
		return "registerTopic"
	case 7:
		return "unregisterTopic"
	default:
		return "unknown"
	}
}
