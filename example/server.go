package main

// do not move
import _ "github.com/KristinaEtc/slflog"

import (
	"github.com/KristinaEtc/bdmq/handler"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/KristinaEtc/config"
	"github.com/ventu-io/slf"
)

//	"github.com/KristinaEtc/bdmq/handler"
//	"github.com/KristinaEtc/bdmq/transport"

var log = slf.WithContext("server-main.go")

// Global is a struct for config
// One Node
type Global struct {
	MachineID string
	Links     []transport.LinkDescFromJSON
}

var globalOpt = Global{
	MachineID: "kristina-note-test",
	Links: []transport.LinkDescFromJSON{
		transport.LinkDescFromJSON{
			LinkID:  "ID1",
			Address: "localhost:6666",
			Mode:    "server",
			Handler: "testHandler",
		},
	},
}

func main() {

	config.ReadGlobalConfig(&globalOpt, "server-main.go")

	transport.RegisterHandlerFactory("testHandler", handler.HandlerTestFactory{})
	n := transport.NewNode()

	err := n.InitLinkDesc(globalOpt.Links)
	if err != nil {
		log.Errorf("InitLinkDesc err: %s", err.Error())
	}

	//n.AddLinkDesc([]byte(fmt.Sprintf("%v", globalOpt.Node)))
	n.Run()

	/*
		//n.SendMessage("topic1", "my message")
		//n.Wait()
		//n.Stop()
	*/
}
