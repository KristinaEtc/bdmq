package stomp

import (
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("stomp.go")

type StompProcesser struct {
	node *StompNode
}

func (sP *StompProcesser) ProcessCommand(cmd transport.Command) (known bool, isExiting bool) {
	var id = cmd.GetCommandID()
	log.Warnf("[stomp] process command=%+v, cmd_id=%s", cmd, sP.CommandToString(cmd.GetCommandID()))

	switch id {
	case stompSendFrameCommand:
		{
			cmdSendFrame, ok := cmd.(*StompCommandSendFrame)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}
			log.Debugf("CommandToString: [%s]; processing: %v", sP.CommandToString(stompSendFrameCommand), cmdSendFrame)

			for y, _ := range sP.node.LinkActives {
				log.Errorf("LinkActives ID=[%s]", y)
			}

			lActive, ok := sP.node.LinkActives[cmdSendFrame.linkActiveID]
			if !ok {
				log.Errorf("Wrong Link Active ID: [%s]; ignored.", cmdSendFrame.linkActiveID)
				return true, false
			} else {
				log.Infof("lActive [%s] is ok", lActive.LinkActiveID)
			}

			frameInByte := lActive.FrameProcessor.ToByte(cmdSendFrame.frame)
			lActive.SendMessageActive(frameInByte)

			return true, false
		}

	default:
		{
			return false, false
		}
	}
}

func (sP *StompProcesser) CommandToString(c transport.CommandID) string {
	switch c {
	case 100:
		return "stompSendFrame"
	default:
		return "unknown"
	}
}
