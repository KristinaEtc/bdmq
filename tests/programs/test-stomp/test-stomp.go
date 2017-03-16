package main

import _ "github.com/KristinaEtc/slflog"
import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/stomp"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/KristinaEtc/config"
	"github.com/ventu-io/slf"
)

// do not move

var log = slf.WithContext("main.go")

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

	config.ReadGlobalConfig(&globalOpt, "stomp.go")
	log.Debugf("config=%v", globalOpt.Links)

	transport.RegisterHandlerFactory("stompHandler", stomp.HandlerStompFactory{})
	transport.RegisterFrameProcessorFactory("stomp", frame.StompFrameProcessorFactory{})
	n := stomp.NewNode()

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
			message := "message:" + "yi"
			frame := frame.New(
				"COMMIT",
				"test:tttt", "test2:2", "test1:1", message,
			)

			n.SendFrame("ID2", frame)
			time.Sleep(time.Second * 3)
			break
			//	n.SendMessage("notProcessedLinkID", scanner.Text()+"\n")
		}
	}

	for scanner.Scan() {
		scanner.Text()
	}
}
