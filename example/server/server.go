package main

// do not move
import _ "github.com/KristinaEtc/slflog"
import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/KristinaEtc/bdmq/tcprec"
	"github.com/KristinaEtc/config"

	"github.com/ventu-io/slf"
)

var log = slf.WithContext("server-main.go")

//Server is a struct for config
type Global struct {
	MachineID string
	Links     []tcprec.LinkOpts
	//Links []Link
	//CallerInfo bool
}

/*
type Link struct {
	TransportProtocol   string
	ApplicationProtocol string
	NodeConf            transport.Noder
}
*/

var globalOpt = Global{
	MachineID: "kristna-note-test",
	Links: []tcprec.LinkOpts{
		//Link{NodeConf: &tcprec.TCPReconn{
		tcprec.LinkOpts{
			ID:         "server1",
			Address:    "localhost:6666",
			Mode:       "server",
			Internal:   5,
			MaxRetries: 10,
		},
	},
	//TransportProtocol:   "UDP",
	//ApplicationProtocol: "STOMP",

}

func run(wg *sync.WaitGroup, done chan (bool), id string, ln net.Listener) {
	/*	defer func() {
			wg.Done()
			done <- true
		}()

		var conn net.Conn
		conn, err := ln.Accept()
		if err != nil {
			log.WithField("id", id).Error(err.Error())
			return
		} else {
			log.WithField("id", id).Info("accepted")
		}

		var buf = make([]byte, 256)

		for {
			_, err := conn.Read(buf)
			switch err {
			case io.EOF:
				log.Debugf("%s detected closed LAN connection", id)
				conn.Close()
				conn = nil
			case tcprec.ErrMaxRetries:
				log.Warn("sever gave up, reached retry limit")
				return
			case nil:
				log.WithField("id", id).Infof("Msg=%s", string(buf))

				answer := strings.ToUpper(string(buf))
				_, err := conn.Write([]byte(answer))
				switch err {
				case io.EOF:
					log.Debugf("%s detected closed LAN connection", id)
					conn.Close()
					conn = nil
					return
				case tcprec.ErrMaxRetries:
					log.Warn("sever gave up, reached retry limit")
					return
				}

			default:
				log.Debug("default")
				time.Sleep(time.Second)
				continue
			}
		}*/
}

func main() {

	/*
		done := make(chan bool)
	*/
	var wg sync.WaitGroup

	// listen on all interfaces
	config.ReadGlobalConfig(&globalOpt, "server-example")
	log.Infof("server configuration: %v\n", globalOpt)

	//
	//log.Debugf("num of nodes=%d", len(globalOpt.Links))

	/*var ln net.Listener
	var err error
	for name, link := range globalOpt.Links {
		log.WithField("id", name).Debug("Running server...")

		 listen on all interfaces
		ln, err = net.Listen("tcp", link.Address)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		go run(&wg, done, link.ID, ln)

			links, numOfFails := transport.ParseConfig(globalOpt.Links)
	if links == nil
	*/

	l, err := tcprec.Init(globalOpt.Links)
	if err != nil {
		log.Errorf("tcprec Init error: %s", err.Error())
		return
	}

	for _, conn := range l {
		wg.Add(1)
		log.Debug("Added 1 waitg")
		go func() {
			// run loop forever (or until ctrl-c)
			for {
				log.Debug("Here")
				// will listen for message to process ending in newline (\n)
				message, _ := bufio.NewReader(conn).ReadString('\n')
				// output message received
				fmt.Print("Message Received:", string(message))
				// sample process for string received
				newmessage := strings.ToUpper(message)
				// send new string back to client
				conn.Write([]byte(newmessage + "\n"))
				log.Debug("Here2")
			}
		}()

	}

	wg.Wait()
}
