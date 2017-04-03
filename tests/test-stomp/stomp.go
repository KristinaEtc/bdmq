package main

import _ "github.com/KristinaEtc/slflog"

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/stomp"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/KristinaEtc/config"

	"github.com/ventu-io/slf"
)

var log = slf.WithContext("main")

var frameReceived int

// Global used for all configs.
type Global struct {
	MachineID      string
	Links          []transport.LinkDescFromJSON
	FileWithFrames string
	ShowFrames     bool
}

var globalOpt = Global{
	MachineID:      "kristina-note-test",
	ShowFrames:     true,
	FileWithFrames: "",
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

func read(ch chan frame.Frame) {
	//	defer func() {
	//		wg.Done()
	//	}()

	for {
		select {
		case _ = <-ch:
			//	if globalOpt.ShowFrames {
			//	log.Infof("Received: %s", fr.Dump())
			//	}
			frameReceived++
		}
	}
}

func write(n *stomp.NodeStomp) error {

	if globalOpt.FileWithFrames == "" {
		log.Errorf("Wrong name of file with frames. Exit.")
		return errors.New("Set a file with frames in command line's flag. Exit.")
	}
	fd, err := os.Open(globalOpt.FileWithFrames)
	if err != nil {
		log.Errorf("Error: %s", err.Error())
		return err
	}

	defer func() {
		fd.Close()
	}()

	reader := frame.NewReader(fd)
	for {
		frame, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				log.Errorf("read from file failed: eof")
			} else {
				log.Errorf("read from file failed: %s", err.Error())
			}
			return err
		}
		if frame == nil {
			//log.Infof("heartbeat")
			continue
		}
		//log.Infof("Sending: [%s], [%v], [%s]", frame.Command, frame.Header, string(frame.Body))

		n.SendFrame("test-topic", *frame)
	}

}

func main() {

	config.ReadGlobalConfig(&globalOpt, "stomp.go")
	log.Infof("FileWithFrames=%s", globalOpt.FileWithFrames)
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

	time.Sleep(time.Second * 5)

	go read(ch)
	write(n)
	time.Sleep(time.Second * 5)
	log.Infof("=============================================1=Frames received: %d", frameReceived)
	frameReceived = 0
	write(n)
	//if err != nil {
	//log.Errorf("Error write: %s", err.Error())
	//n.Stop()
	//}
	time.Sleep(time.Second * 5)
	log.Infof("=============================================2=Frames received: %d", frameReceived)
}
