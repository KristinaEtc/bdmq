package main

// do not move
import _ "github.com/KristinaEtc/slflog"

import (
	"fmt"

	"github.com/KristinaEtc/bdmq/tcprec"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("server-main.go")

//Global is a struct for config
type Global struct {
	MachineID     string
	MaxRetryCount int
	Links         []tcprec.LinkOpt
	//Links []Link
	//CallerInfo bool
}

var globalOpt = Global{
	MachineID: "kristina-note-test",
	Nodes: {
		[]tcprec.LinkOpts{
			//Link{NodeConf: &tcprec.TCPReconn{
			tcprec.LinkOpts{
				ID:       "server1",
				Address:  "localhost:6666",
				Mode:     "server",
				Internal: 5,
				Handler:  "stomp",
			},
		},
		//
		//TransportProtocol:   "UDP",
		//ApplicationProtocol: "STOMP",
	},
}

func main() {

	n := tcprec.NewNode()
	err := n.InitLinkDesc([]byte(fmt.Sprintf("%v", globalOpt.Node)))
	if err != nil {
		log.Errorf("InitLinkDesc err: %s", err.Error())
	}
	//n.AddLinkDesc([]byte(fmt.Sprintf("%v", globalOpt.Node)))
	n.Run()

	//n.SendMessage("topic1", "my message")
	//n.Wait()
	//n.Stop()
}
