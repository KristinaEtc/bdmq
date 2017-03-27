package main

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/KristinaEtc/bdmq/stomp"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/KristinaEtc/config"
	_ "github.com/KristinaEtc/slflog"
	"github.com/go-stomp/stomp/frame"
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

/*
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
				"topic:topic-test", "test2:2", "test1:1", message,
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

func subscribe() (map[string]chan[]byte, bool) {

	var subscriptions = make(mapp[string]chan[]byte)
	for _, link := range globalOpt.Links {
		if len(link.Topics) != 0 {
			for _, topic := range link.Topics {
				ch, err = n.Subscribe(link.LinkID, topic)
				if err != nil{
					log.Errorf("could not subscribe link [%s] for topic [%s]: %s", link.LinkID, topic, err.Error())
					return nil, err
				}
				subscriptions[topic] = ch
			}
		}
	}
	return false
}
*/

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
				"from:ID-2", "topic:test-topic", "test1:1", message,
			)

			n.Subscribe("test-topic")
			n.SendFrame("ID-2", frame)
			time.Sleep(time.Second * 3)
			break
			//	n.SendMessage("notProcessedLinkID", scanner.Text()+"\n")
		}
	}

	for scanner.Scan() {
		scanner.Text()
	}
}
