package main

// do not move
import (
	"bufio"
	"fmt"

	"github.com/KristinaEtc/bdmq/tcprec"
	"github.com/KristinaEtc/config"
	_ "github.com/KristinaEtc/slflog"

	_ "github.com/KristinaEtc/slflog"
	"github.com/ventu-io/slf"
)

//"github.com/KristinaEtc/config"

var log = slf.WithContext("client.go")

//To test the library, you can run a local TCP server with:
//$ ncat -l 9999 -k

//Server is a struct for config
type Global struct {
	MachineID string
	Links     []tcprec.LinkOpts
	//CallerInfo bool
}

var globalOpt = Global{
	MachineID: "kristina-note-clients",
	Links: []tcprec.LinkOpts{
		tcprec.LinkOpts{
			ID:         "user1",
			Address:    "localhost:1234",
			Mode:       "client",
			Internal:   5,
			MaxRetries: 10,
		},
		//tcprec.LinkOpts{
		//	ID:         "user2",
		//	Address:    "localhost:7777",
		//	Mode:       "client",
		//	Internal:   2,
		//	MaxRetries: 7,
		//},
	},
	//	CallerInfo: false,
}

func main() {

	config.ReadGlobalConfig(&globalOpt, "client")
	log.Infof("main: %v\n", globalOpt)

	conns, err := tcprec.Init(globalOpt.Links)
	if err != nil {
		log.Error(err.Error())
	}

	//if tcprec.LinkExists(serverID) {
	for name, coo := range conns {
		log.Infof("name=%s", name)
		if coo == nil {
			log.Error("tttggggggrrr")
		}
	}

	//var serverID
	var msg string
	gg := conns["kristina-client"]

	for {
		//fmt.Print("node ID> ")
		//fmt.Scanf("node ID> %s\n", &serverID)
		//fmt.Scanln(&serverID)
		//fmt.Println(serverID)
		fmt.Print("message> ")
		fmt.Scanln(&msg)
		if msg == "/q" {
			fmt.Println("goodbye")
			break
		}

		gg.Write([]byte(msg + "\n"))

		message, _ := bufio.NewReader(gg).ReadString('\n')
		// output message received
		fmt.Print("Message Received:", string(message))

	}
}
