package main

import (
	"fmt"

	_ "github.com/KristinaEtc/slflog"

	"github.com/KristinaEtc/bdmq/transport/tcp"
	"github.com/KristinaEtc/config"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("main.go")

var (
	// These fields are populated by govvv
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
	Version    string
)

type ClientConf struct {
	BindingAddr         string
	TransportProtocol   string
	ApplicationProtocol string
}

type NodesConf struct {
	Clients []ClientConf
}

// ConfFile is a file with all program options
type ConfFile struct {
	Name string
	Node NodesConf
}

var globalOpt = ConfFile{
	Name: "init-system-k",
	Node: NodesConf{
		Clients: []ClientConf{},
	},
}

func main() {

	config.ReadGlobalConfig(&globalOpt, "node options")

	log.Infof("%s", globalOpt.Name)

	log.Error("----------------------------------------------")

	log.Infof("BuildDate=%s\n", BuildDate)
	log.Infof("GitCommit=%s\n", GitCommit)
	log.Infof("GitBranch=%s\n", GitBranch)
	log.Infof("GitState=%s\n", GitState)
	log.Infof("GitSummary=%s\n", GitSummary)
	log.Infof("VERSION=%s\n", Version)

	log.Info("Starting working...")

	tcpConnFactory := &tcp.TCPNodeFactory{}

	node := tcpConnFactory.Init(tcpConnFactory,
		[]byte(fmt.Sprintf("%v", globalOpt.Node)),
		"Handler, which I didn't implement")

	node.Write("yo")

}
