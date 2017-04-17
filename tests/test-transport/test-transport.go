package main

import _ "github.com/KristinaEtc/slflog"

import (
	"bufio"
	"os"
	"strings"

	"github.com/KristinaEtc/bdmq/handlers"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/KristinaEtc/config"

	"github.com/ventu-io/slf"
)

var log = slf.WithContext("main.go")

// Global is a struct with all configs; used with github.com/KristinaEtc/config package.
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

	config.ReadGlobalConfig(&globalOpt, "main.go")
	log.Debugf("config=%v", globalOpt.Links)

	transport.RegisterHandlerFactory("echoHandler", handlers.HandlerEchoFactory{})
	transport.RegisterHandlerFactory("testHandler", handlers.HandlerTestFactory{})

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
}