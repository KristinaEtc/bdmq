package transport

//type Command interface{}

type CommandProcesser interface {
	ProcessCommand(*NodeCommand) bool
}

type DefaultProcesser struct {
	node *Node
}

func (dP *DefaultProcesser) ProcessCommand(cmdMsg *NodeCommand) (isExiting bool) {
	log.Debugf("process command=%s", cmdMsg.cmd.String())

	n := dP.node

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
