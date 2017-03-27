package main

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/stomp"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/KristinaEtc/config"
	_ "github.com/KristinaEtc/slflog"
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

func read(wg *sync.WaitGroup, ch chan frame.Frame) {
	defer func() {
		wg.Done()
	}()

	for {
		select {
		case fr := <-ch:
			log.Infof("Got frame: %s", fr.Dump())
		}
	}
}

func write(wg *sync.WaitGroup, n *stomp.NodeStomp) {
	defer func() {
		wg.Done()
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if strings.ToLower(scanner.Text()) == "q" {
			n.Stop()
			break
		} else {

			message := scanner.Text() + "yi"
			frame := frame.New(
				"COMMIT",
				"from:ID-2", "topic:test-topic", "test1:1", message,
			)

			n.SendFrame("test-topic", *frame)

			//	n.SendMessage("notProcessedLinkID", scanner.Text()+"\n")
		}
	}
}

func main() {

	config.ReadGlobalConfig(&globalOpt, "stomp.go")
	log.Debugf("config=%+v", globalOpt.Links)

	transport.RegisterHandlerFactory("stomp", stomp.HandlerStompFactory{})

	n := stomp.NewNode()
	n.AddCmdProcessor(&stomp.ProcessorStomp{Node: n})

	err := n.InitLinkDesc(globalOpt.Links)
	if err != nil {
		log.Errorf("InitLinkDesc error: %s", err.Error())
		os.Exit(0)
	}

	err = n.Run()
	if err != nil {
		log.Errorf("Run error: %s", err.Error())
	}

	ch, err := n.Subscribe("test-topic")
	if err != nil {
		log.Errorf("Could not subscribe: %s", err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go read(&wg, ch)
	go write(&wg, n)

	wg.Wait()

}
