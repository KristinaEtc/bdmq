package main

// do not move
import _ "github.com/KristinaEtc/slflog"
import (
	"io/ioutil"
	"strconv"
	"time"

	"github.com/KristinaEtc/bdmq/tcprec"
	"github.com/KristinaEtc/config"

	"github.com/ventu-io/slf"
)

//"github.com/KristinaEtc/config"

var log = slf.WithContext("server.go")

//To test the library, you can run a local TCP server with:
//$ ncat -l 9999 -k

//Server is a struct for config
type Global struct {
	Addr string
}

type Link struct{
	ID string
	Address string
	Mode string
	Internal int
}

var globalOpt = Global{
	Links: []Link{
		Link{
			ID: "user",
			Address: "localhost:6666",
			Mode: "client",
			Internal:5,
		}
	}
}

func main() {

	config.ReadGlobalConfig(&globalOpt, "client-example")
	var conns = tcprec.Init(globalOpt)

	// connect to a TCP server
/*	conn, err := tcprec.ConnectClient("tcp", globalOpt.Addr)
	if err != nil {
		log.Errorf("%s", err.Error())
		return
	}*/

	

	//	conn.SetMaxRetries(5)

	//buf := make([]byte, 256)

	//go send(conn)
	//go recv(conn)
	// client sends "hello, world!" to the server every second
	for i := 0; i < 10; i++ {
		_, err := conn.Write([]byte("hello, world!" + strconv.Itoa(i)))
		if err != nil {
			log.Error(err.Error())
			// if the client reached its retry limit, give up
			if err == tcprec.ErrMaxRetries {
				log.Warn("client gave up, reached retry limit")
				return
			}
		}
		time.Sleep(time.Second)
		result, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Error(err.Error())
			// if the client reached its retry limit, give up
			//if err == tcprec.ErrMaxRetries {
			//	log.Warn("client gave up, reached retry limit")
			return
		}
		log.Infof(string(result))
	}

	/*log.Debug("test")
	_, err = conn.Read(buf)
	if err != nil {
		//return
		if err == tcprec.ErrMaxRetries {
			log.Warn("client gave up, reached retry limit")
			return
		} else {
			continue
		}
	}*/
	//log.Info(string(buf))

	// not a tcprec error, just panic
	//log.Error(err.Error())
	//}

	/*_, err = conn.Read(buf)
	if err != nil {
		log.Error(err.Error())
	}

	log.Info("client says hello!")
	time.Sleep(time.Second)*/
	//	}
}

/*func send(conn net.Conn) {
	for {
		_, err = conn.Read(buf)
		if err != nil {
			log.Error(err.Error())
		}

		log.Info("client says hello!")
		time.Sleep(time.Second)
	}
}

func recv(conn net.Conn) {
	for {
		_, err := conn.Write([]byte("hello, world!\n"))
		if err != nil {
			// if the client reached its retry limit, give up
			if err == tcprec.ErrMaxRetries {
				log.Warn("client gave up, reached retry limit")
				return
			}
		}
	}
}*/
