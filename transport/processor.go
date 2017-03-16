package transport

type Command interface {
	GetCommandID() CommandID
}

type CommandProcesser interface {
	ProcessCommand(Command) (bool, bool)
	CommandToString(CommandID) string
}

type DefaultProcesser struct {
	node    *Node
	handler Handler
}

func (dP *DefaultProcesser) ProcessCommand(cmd Command) (known bool, isExiting bool) {
	var id = cmd.GetCommandID()
	log.Debugf("process command=%+v, cmd_id=%s", cmd, dP.CommandToString(cmd.GetCommandID()))

	n := dP.node

	switch id {
	case registerActive:
		{
			cmdActive, ok := cmd.(*NodeCommandActiveLink)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}

			n.LinkActives[cmdActive.active.Id()] = cmdActive.active
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
					log.Debugf("%s", lA.Id())
					//return false
				}
			} else {
				log.Debug("NIL")
			}

			delete(n.LinkActives, cmdActive.active.Id())
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

			n.LinkControls[cmdControl.ctrl.getId()] = cmdControl.ctrl
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

			delete(n.LinkControls, cmdControl.ctrl.getId())
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
			return true, false
		}
	case sendMessageNode:
		{
			cmdMessage, ok := cmd.(*NodeCommandSendMessage)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}

			if len(n.LinkActives) > 0 {
				for _, lA := range n.LinkActives {
					log.Debug("for testing i'm choosing all active links")
					lA.SendMessageActive([]byte(cmdMessage.msg))
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

func (dP *DefaultProcesser) CommandToString(c CommandID) string {
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
		return "sendMessage"
	default:
		return "unknown"
	}
}
