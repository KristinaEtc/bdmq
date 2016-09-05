package main

import (
	"log"
	"time"

	"github.com/KristinaEtc/bdmq/tcprec"
)

//To test the library, you can run a local TCP server with:
//$ ncat -l 9999 -k

func main() {
	// connect to a TCP server
	conn, err := tcprec.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Fatal(err)
	}

	// client sends "hello, world!" to the server every second
	for {
		_, err := conn.Write([]byte("hello, world!\n"))
		if err != nil {
			// if the client reached its retry limit, give up
			if err == tcprec.ErrMaxRetries {
				log.Println("client gave up, reached retry limit")
				return
			}
			// not a tcprec error, just panic
			log.Fatal(err)
		}
		log.Println("client says hello!")
		time.Sleep(time.Second)
	}
}
