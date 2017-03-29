package stomp

import (
	"github.com/KristinaEtc/bdmq/frame"
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

			/*
				subcriptionData, ok := s.Node.subscriptions[cmdSendFrame.topic]
				if !ok {
					log.Warnf("Wrong topic name: %s; ignored.", cmdSendFrame.topic)
					return true, false
				}

				lActive := s.Node.LinkActives[subcriptionData.linkActiveID]
			*/
			log.Debug("Now I'm sending a frame to the 1st handler")
			for _, h := range s.Node.handlers {
				h.OnWrite(cmdSendFrame.frame)
				//break
			}

			return true, false
		}

	case stompReceiveFrameCommand:
		{
			cmdReceiveFrame, ok := cmd.(*CommandReceiveFrameStomp)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}
			log.Debugf("Command=[%s/%d]; frame: [%s]", s.CommandToString(stompReceiveFrameCommand), stompReceiveFrameCommand, cmdReceiveFrame.frame.Dump())

			// getting a topic from hiddr and sending to channels which subscribed for them;
			// now not implemented

			subData, ok := s.Node.subscriptions[cmdReceiveFrame.topic]
			if !ok {
				log.Warnf("Got message for topic=[%s]; nobody subscribed for it.", cmdReceiveFrame.topic)
				return true, false
			}
			subData.ch <- cmdReceiveFrame.frame

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
				chTopic := make(chan frame.Frame, 0)
				s.Node.subscriptions[cmdTopic.topic] = Subscription{
					ch: chTopic,
				}
			}

			var res subscripted
			res.ok = true
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
	case stompRegisterStompHandlerCommand:
		{

			cmdHandlerRegister, ok := cmd.(*CommandRegisterHandlerStomp)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}

			s.Node.handlers = append(s.Node.handlers, cmdHandlerRegister.handler)

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
	case 101:
		return "stompSubscribe"
	case 102:
		return "stompUnsubscribe"
	case 103:
		return "stompReceiveFrame"
	case 104:
		return "stompRegisterStompHandlerCommand"
	default:
		return "unknown"
	}
}
