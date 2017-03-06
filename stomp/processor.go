package stomp

import (
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("stomp.go")

type StompProcesser struct {
	node *transport.Node
}

func (sP *StompProcesser) ProcessCommand(cmd transport.Command) (known bool, isExiting bool) {
	var id = cmd.GetCommandID()
	log.Debugf("process command=%+v", cmd)

	switch id {
	case testCommand:
		{
			cmdTest, ok := cmd.(*StompCommandTestLink)
			if !ok {
				log.Errorf("Invalid command type %v", cmd)
				return false, true
			}
			log.Debugf("Successful %s processing: %v", sP.CommandToString(testCommand), cmdTest)
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
		return "testCommand"
	default:
		return "unknown"
	}
}
