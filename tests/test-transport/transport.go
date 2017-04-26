package main

import _ "github.com/KristinaEtc/slflog"
import (
	"os"

	"github.com/KristinaEtc/bdmq/handlers"

	test "github.com/KristinaEtc/bdmq/test-commands"

	test_transport "github.com/KristinaEtc/bdmq/test-commands/transport"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/KristinaEtc/config"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("test-transport")

// Global is a struct with all configs; used with github.com/KristinaEtc/config package.
type Global struct {
	MachineID        string
	Links            []transport.LinkDescFromJSON
	FileWithCommands string
}

var globalOpt = Global{
	MachineID:        "kristina-note-test",
	FileWithCommands: "commands.cmd",
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

	config.ReadGlobalConfig(&globalOpt, "test-transport")
	log.Debugf("config=%v", globalOpt.Links)

	transport.RegisterHandlerFactory("echoHandler", handlers.HandlerEchoFactory{})
	transport.RegisterHandlerFactory("testHandler", handlers.HandlerTestFactory{})
	transport.RegisterHandlerFactory("helloHandler", handlers.HandlerHelloWorldFactory{})

	n := transport.NewNode()
	err := n.InitLinkDesc(globalOpt.Links)
	if err != nil {
		log.Errorf("InitLinkDesc error: %s", err.Error())
		os.Exit(0)
	}

	err = n.Run()
	if err != nil {
		log.Errorf("Run error: %s", err.Error())
	}
	/*
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if strings.ToLower(scanner.Text()) == "q" {
				n.Stop()
				//time.Sleep(time.Second * 5)
				break
				//os.Exit(0)
			} else {
				n.SendMessage("notProcessedLinkID", scanner.Text()+"\n")
			}
		}
	*/

	cmdCtx := test.NewCommandsRegistry()
	test_transport.Register(n, &cmdCtx)

	test.Process(&cmdCtx, globalOpt.FileWithCommands)
}
