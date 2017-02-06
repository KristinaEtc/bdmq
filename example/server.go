package main

// do not move
import _ "github.com/KristinaEtc/slflog"

import (
	"bufio"
	"os"
	"strings"

	"github.com/KristinaEtc/bdmq/handlers"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/KristinaEtc/config"
	_ "github.com/KristinaEtc/slflog"
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
			BufSize: 1024,
		},
	},
}

func main() {

	config.ReadGlobalConfig(&globalOpt, "server-main.go")

	log.Debugf("config=%v", globalOpt.Links)

	transport.RegisterHandlerFactory("testHandler", handlers.HandlerTestFactory{})
	transport.RegisterHandlerFactory("echoHandler", handlers.HandlerEchoFactory{})
	transport.RegisterHandlerFactory("worldHandler", handlers.HandlerHelloWorldFactory{})

	n := transport.NewNode()

	err := n.InitLinkDesc(globalOpt.Links)
	if err != nil {
		log.Errorf("InitLinkDesc err: %s", err.Error())
	}

	//n.AddLinkDesc([]byte(fmt.Sprintf("%v", globalOpt.Node)))
	n.Run()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		if strings.ToLower(scanner.Text()) == "q" {
			n.Stop()
			os.Exit(1)
		}
	}

	//n.SendMessage("topic1", "my message")

}
