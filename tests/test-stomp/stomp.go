package main

import _ "github.com/KristinaEtc/slflog"

import (
	"io"
	"os"

	"github.com/KristinaEtc/bdmq/frame"
	"github.com/KristinaEtc/bdmq/stomp"
	test "github.com/KristinaEtc/bdmq/test-commands"
	test_stomp "github.com/KristinaEtc/bdmq/test-commands/stomp"
	"github.com/KristinaEtc/bdmq/transport"
	"github.com/KristinaEtc/config"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("test-stomp")

// Global used for all configs.
type Global struct {
	MachineID        string
	Links            []transport.LinkDescFromJSON
	FileWithFrames   string
	FileWithCommands string
	ShowFrames       bool
	StopTimeout      int
}

var globalOpt = Global{
	MachineID:        "test",
	ShowFrames:       true,
	FileWithFrames:   "",
	FileWithCommands: "commands.cmd",
	StopTimeout:      5,
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

func write(n *stomp.NodeStomp) error {

	fd, err := os.Open(globalOpt.FileWithFrames)
	if err != nil {
		log.Errorf("File with frames=[%s]: %s", globalOpt.FileWithFrames, err.Error())
		return err
	}

	defer func() {
		if err := fd.Close(); err != nil {
			log.Errorf("File with frames=[%s]: %s", globalOpt.FileWithFrames, err.Error())
		}
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

	config.ReadGlobalConfig(&globalOpt, "test-transport")
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

	defer func() {
		err := n.Stop(globalOpt.StopTimeout)
		if err != nil {
			os.Exit(1)
		}
	}()
	/*
		ch, err := n.Subscribe("test-topic")
		if err != nil {
			log.Errorf("Could not subscribe: %s", err.Error())
			return
		}

		time.Sleep(time.Second * 5)


		go read(ch)
		write(n)
	*/
	cmdCtx := test.NewCommandsRegistry()

	//	test_transport.Register(n.Node, &cmdCtx)
	test_stomp.Register(n, &cmdCtx, true)

	test.Process(&cmdCtx, globalOpt.FileWithCommands)
}
