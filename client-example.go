package main

// do not move
import _ "github.com/KristinaEtc/slflog"

import (
	"time"

	"github.com/KristinaEtc/bdmq/tcprec"
	//"github.com/KristinaEtc/config"

	"github.com/ventu-io/slf"
)

var log = slf.WithContext("server.go")

//To test the library, you can run a local TCP server with:
//$ ncat -l 9999 -k

//Server is a struct for config
type Global struct {
	Addr string
}

var globalOpt = Global{
	Addr: "localhost:7777",
}

func main() {

	//config.ReadGlobalConfig(&globalOpt, "client-example")

	// connect to a TCP server
	conn, err := tcprec.Dial("tcp", globalOpt.Addr)
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}

	buf := make([]byte, 256)
	// client sends "hello, world!" to the server every second
	for {
		/*	_, err := conn.Write([]byte("hello, world!\n"))
			if err != nil {
				// if the client reached its retry limit, give up
				if err == tcprec.ErrMaxRetries {
					log.Warn("client gave up, reached retry limit")
					return
				}

				// not a tcprec error, just panic
				log.Error(err.Error())
			}*/

		_, err = conn.Read(buf)
		if err != nil {
			log.Error(err.Error())
		}

		log.Info("client says hello!")
		time.Sleep(time.Second)
	}
}
