package stomp

import "github.com/KristinaEtc/bdmq/transport"

type StompCommandTestLink struct {
	transport.NodeCommand
	message string
}

const (
	testCommand transport.CommandID = 100
)
