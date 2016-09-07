package main

// do not move
import (
	"net"
	"strings"
	"time"

	_ "github.com/KristinaEtc/slflog"
	"github.com/ventu-io/slf"
)

var log = slf.WithContext("server.go")

func main() {

	log.Info("Launching server...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", ":7777")
	if err != nil {
		log.Error(err.Error())
		return
	}
	// accept connection on port
	conn, err := ln.Accept()
	if err != nil {
		log.Error(err.Error())
		return
	} else {
		log.Info("accepted")
	}

	var buf = make([]byte, 256)

	// run loop forever (or until ctrl-c)
	for {
		// will listen for message to process ending in newline (\n)
		/*reader := bufio.NewReader(conn)
		/*	if err != nil {
			log.Error(err.Error())
			return
		}

		_, err = reader.Read(buf)
		if err != nil {
			log.Error(err.Error())
			return
		}*/

		conn.Read(buf)
		// output message received
		log.Infof("Message Received: %s\n", string(buf))
		// sample process for string received
		newmessage := strings.ToUpper(string(buf))
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))

		time.Sleep(time.Second * 1)

		//conn.Write([]byte(strconv.Itoa(i) + "yo\n"))
		/*writer := bufio.NewWriter(conn)
		_, err := writer.WriteString(strconv.Itoa(i) + "yo\n")
		if err != nil {
			log.Error(err.Error())
			return
		}*/
		//time.Sleep(time.Second)
	}
}

/*
package main

import (
	"github.com/KristinaEtc/bdmq/tcprec"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// open a server socket
	s, err := net.Listen("tcp", "localhost:7777")
	if err != nil {
		log.Fatal(err)
	}
	// save the original port
	addr := s.Addr()

	// connect a client to the server
	c, err := tcprec.Dial("tcp", s.Addr().String())
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// shut down and boot up the server randomly
	var swg sync.WaitGroup
	swg.Add(1)
	go func() {
		defer swg.Done()
		for i := 0; i < 5; i++ {
			log.Println("server up")
			time.Sleep(time.Millisecond * 100 * time.Duration(rand.Intn(20)))
			if err := s.Close(); err != nil {
				log.Fatal(err)
			}
			log.Println("server down")
			time.Sleep(time.Millisecond * 100 * time.Duration(rand.Intn(20)))
			s, err = net.Listen("tcp", addr.String())
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	// client writes to the server and reconnects when it has to
	// this is the interesting part
	var cwg sync.WaitGroup
	cwg.Add(1)
	go func() {
		defer cwg.Done()
		for {
			if _, err := c.Write([]byte("hello, world!\n")); err != nil {
				switch e := err.(type) {
				case tcprec.Error:
					if e == tcprec.ErrMaxRetries {
						log.Println("client leaving, reached retry limit")
						return
					}
				default:
					log.Fatal(err)
				}
			}
			log.Println("client says hello!")
		}
	}()

	// terminates the server indefinitely
	swg.Wait()
	if err := s.Close(); err != nil {
		log.Fatal(err)
	}

	// wait for the client to give up
	cwg.Wait()
}
*/
