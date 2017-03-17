package stomp

import (
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("NodeStomp")

// ProcessorStomp inherited from transport.Processor
type ProcessorStomp struct {
	node *NodeStomp
}

// ProcessCommand process STOMP commands.
func (s *ProcessorStomp) ProcessCommand(cmd transport.Command) (known bool, isExiting bool) {
	var id = cmd.GetCommandID()
	log.Debugf("Process command [%s]", s.CommandToString(cmd.GetCommandID()))

	switch id {
	case stompSendFrameCommand:
		{
			cmdSendFrame, ok := cmd.(*CommandSendFrameStomp)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}
			log.Debugf("Command=[%s/%d]; frame: [%s]", s.CommandToString(stompSendFrameCommand), stompSendFrameCommand, cmdSendFrame.frame.Dump())

			lActives, ok := s.node.Topics[cmdSendFrame.topic]
			if !ok {
				log.Warnf("Wrong topic name: %s; ignored.", cmdSendFrame.topic)
				return true, false
			}

			for _, lA := range lActives {
				frameInByte := lA.FrameProcessor.ToByte(cmdSendFrame.frame)
				lA.SendMessageActive(frameInByte)
			}

			return true, false
		}

	default:
		{
			return false, false
		}
	}
}

// CommandToString returns a string representation of command's ID.
// CommandToString is a methor of transport.CommandProcessor.
func (s *ProcessorStomp) CommandToString(c transport.CommandID) string {
	switch c {
	case 100:
		return "stompSendFrame"
	default:
		return "unknown"
	}
}
