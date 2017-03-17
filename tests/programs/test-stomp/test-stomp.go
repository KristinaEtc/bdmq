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

var log = slf.WithContext("main")

// Global used for all configs.
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

func read(done chan bool, n *stomp.NodeStomp) {
	defer func() {
		done <- true
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if strings.ToLower(scanner.Text()) == "q" {
			n.Stop()
			break
		} else {
			message := "message:" + scanner.Text()
			frame := frame.New(
				"COMMIT",
				"test:tttt", "test2:2", "test1:1", message,
			)

			n.SendFrame("topic-test", frame)
		}
	}
}

func listen(done chan bool, n *stomp.NodeStomp) {
	defer func() {
		done <- true
	}()

	var ch chan []byte
	var err error
	for {
		ch, err = n.GetChannel("ID2")
		if err != nil {
			log.Errorf("Could not get link channel: %s", err.Error())
			time.Sleep(time.Second * 1)
			continue
		}
		log.Info("Got channel")
		break
	}

	select {
	case msgByte := <-ch:
		log.Infof("Got frame: %s", string(msgByte))
	}
}

func main() {

	config.ReadGlobalConfig(&globalOpt, "stomp.go")
	log.Debugf("config=%+v", globalOpt.Links)

	transport.RegisterHandlerFactory("stomp", stomp.HandlerStompFactory{})
	transport.RegisterFrameProcessorFactory("stomp", frame.FactoryStomp{})
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

	done := make(chan bool)
	go listen(done, n)
	go read(done, n)

	<-done
}
