package stomp

import (
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("NodeStomp")

// ProcessorStomp inherited from transport.Processor
type ProcessorStomp struct {
	Node *NodeStomp
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

			lActive, ok := s.Node.LinkActives[cmdSendFrame.linkActiveID]
			if !ok {
				log.Warnf("Wrong linkActiveID name: %s; ignored.", cmdSendFrame.linkActiveID)
				return true, false
			}

			lActive.GetHandler().(*HandlerStomp).OnWrite(cmdSendFrame.frame)

			return true, false
		}

	case stompSubscribeCommand:
		{
			cmdTopic, ok := cmd.(*CommandSubscribeStomp)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}

			log.Infof("topic=%v", cmdTopic.topic)

			_, ok = s.Node.subscriptions[cmdTopic.topic]
			if !ok {
				s.Node.subscriptions[cmdTopic.topic] = make(map[string]*chan transport.Frame, 0)
			}

			var res subscripted

			ch, err := s.Node.LinkActives[cmdTopic.linkActiveID].GetHandler().(*HandlerStomp).Subscribe(cmdTopic.topic)
			if err != nil {
				res.ok = false
				res.err = err
				cmdTopic.sub <- res
				return true, false // isExiting?
			}
			res.ok = true
			res.err = nil
			s.Node.subscriptions[cmdTopic.topic][cmdTopic.linkActiveID] = ch
			cmdTopic.sub <- res
			return true, false

		}

		/*
			case unregisterTopic:
				{
					cmdTopic, ok := cmd.(*NodeCommandTopic)
					if !ok {
						log.Errorf("Invalid command type %v", cmd)
						return false, true
					}

					topicMap, ok := n.Topics[cmdTopic.topicName]
					if !ok {
						log.Warnf("No such topic %s; ignored.", cmdTopic.topicName)
						return true, false
					}

					_, ok = topicMap[cmdTopic.active.ID()]
					if !ok {
						log.Warnf("LinkActive with id=%s doesn't subscribe for topic %s; ignored.", cmdTopic.active.ID(), cmdTopic.topicName)
						return true, false
					}

					delete(topicMap, cmdTopic.active.ID())
					if len(topicMap) == 0 {
						delete(n.Topics, cmdTopic.topicName)
					}

					log.Debugf("[unregisterTopic] Lin kActive with id=%s from topic=%d", cmdTopic.active.ID(), cmdTopic.topicName)
					return true, false
				}
		*/
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
	case 101:
		return "stompSubscribe"
	case 102:
		return "stompUnsubscribe"
	default:
		return "unknown"
	}
}
