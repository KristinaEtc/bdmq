package transport

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

// LinkActiver is an interface for connection struct.
// Realization of this interface depends of connection state:
// * connected
// * reconnecting
// * error state
/*
type LinkActiver interface {
	Write(string) error
	Read() error
	GetLinkActiveID() string
	Disconnect() error
}
*/

// LinkControlServer initialize LinkStater interface when
// connecting or reconnecting
type LinkControlServer struct {
	CommandCh chan string
	Node      *Node
	//Exiting      bool
	LinkDesc *LinkDesc
}

func (lS *LinkControlServer) Write(msg string) error {
	log.Warn("(lS *LinkControlServer) Write")
	return fmt.Errorf("%s", ErrStubFunction)
}

func (lS *LinkControlServer) Read() error {
	log.Warn("(lS *LinkControlServer) Read")
	return fmt.Errorf("%s", ErrStubFunction)
}

func (lc *LinkControlServer) Listen() (net.Listener, error) {
	log.Warn("(lc *LinkControlServer) Listen")

	var err error
	var ln net.Listener

	var secToRecon = time.Duration(time.Second * 2)
	var numOfRecon = 0

	for {
		ln, err = net.Listen(network, lc.LinkDesc.address)
		if err == nil {
			log.WithField("ID=", lc.LinkDesc.linkID).Debugf("Created listener with: %s", ln.Addr().String())
			return ln, nil
		}

		log.WithField("ID=", lc.LinkDesc.linkID).Errorf("Error listen: %s. Reconnecting after %d milliseconds.", err.Error(), secToRecon/1000000.0)
		ticker := time.NewTicker(secToRecon)

		select {
		case _ = <-ticker.C:
			{
				if secToRecon < backOffLimit {

					randomAdd := secToRecon / 100 * (20 + time.Duration(r1.Int31n(10)))
					secToRecon = secToRecon*2 + time.Duration(randomAdd)
					numOfRecon++
				}
				ticker = time.NewTicker(secToRecon)
				continue
			}

		case command := <-lc.CommandCh:
			{
				log.Debugf("Listen command %s", command)
				if command == "quit" {
					log.WithField("ID=", lc.LinkDesc.linkID).Info("Got quit command. Closing link.")
					return nil, fmt.Errorf("quit requested")
				}
				log.WithField("ID=", lc.LinkDesc.linkID).Warnf("Got impermissible command %s. Ignored.", command)
			}
		}
	}

	//return fmt.Errorf("%s", ErrStubFunction)
}

func (lc *LinkControlServer) Accept(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.WithField("ID=", lc.LinkDesc.linkID).Errorf("Error accept: %s", err.Error())
			lc.NotifyError(err)
			return
		}
		log.WithField("ID=", lc.LinkDesc.linkID).Debug("New client")

		//n.InitLinkActive(linkD, &conn, ((n.LinkStubs[linkD.linkID]).(*LinkStub)).commandCh)
		conn.Close()
	}
}

func (lc *LinkControlServer) WaitCommand(listener net.Listener) bool {
	select {
	case command := <-lc.CommandCh:
		if command == "quit" {
			log.Debugf("quit received")
			listener.Close()
			return true
		}
		return false
	}
}

func (lS *LinkControlServer) Disconnect() error {
	log.Info("(lS *LinkControlServer)  Disconnect")
	return fmt.Errorf("%s", ErrStubFunction)
}

func (lc *LinkControlServer) Close() {
	log.Info("(lS *LinkControlServer)  Close")
	lc.CommandCh <- "quit"
}

func (lc *LinkControlServer) NotifyError(err error) {
	lc.CommandCh <- "error " + err.Error()
}

func (lc *LinkControlServer) Id() string {
	return lc.LinkDesc.linkID
}

// LinkControlClient initialize LinkStater interface when
// connecting or reconnecting
type LinkControlClient struct {
	CommandCh    *chan string
	LinkActiveID string
}

func (lC *LinkControlClient) Close() error {
	log.Warn("(lC *LinkControlClient) Close")
	return fmt.Errorf("%s", ErrStubFunction)
}

/*func (lC *LinkControlClient) Write(msg string) error {
	log.Warn("(lC *LinkControlClient) Write")
	return fmt.Errorf("%s", ErrStubFunction)
}*/

func (lC *LinkControlClient) GetLinkActiveID() string {
	return lC.LinkActiveID
}

func (lC *LinkControlClient) Read() error {
	log.Warn("(lC *LinkControlClient) Read")
	return fmt.Errorf("%s", ErrStubFunction)
}

func (lC *LinkControlClient) Disconnect() error {
	log.Info("(lC *LinkControlClient)  Disconnect")
	return fmt.Errorf("%s", ErrStubFunction)
}

// LinkActive initialize LinkStater interface when
// connecting is ok
type LinkActive struct {
	conn         *net.Conn
	lock         sync.RWMutex
	status       int32
	LinkConf     *LinkDesc
	LinkActiveID string
	handler      *Handler
	//parser       *parser.Parser
	commandCh *chan string
}

func (lA *LinkActive) GetLinkActiveID() string {
	return lA.LinkActiveID
}

func (lA *LinkActive) Disconnect() error {
	log.Info("(lSub *LinkActive)Disconnect")
	return nil
}

func (lA *LinkActive) Write(msg string) error {

	log.Info("(lSub *LinkActive)Write")
	(*lA.conn).Write([]byte(msg + "\n"))
	return nil
}

//func (lA *LinkActive) Read(b []byte) error {
func (lA *LinkActive) Read() error {
	log.Info("(lSub *LinkActive)Read")

	buf := make([]byte, SizeOfBuf)

	for {
		message, err := bufio.NewReader(*lA.conn).Read(buf)
		if err != nil {
			log.WithField("ID=", lA.LinkActiveID).Errorf("Error read: %s", err.Error())
			log.WithField("ID=", lA.LinkActiveID).Errorf("ADD GRASEFUL RECONNECT")
			os.Exit(1)
		}
		log.WithField("ID=", lA.LinkActiveID).Debugf("Message Received: %s", string(message))
		go (*lA.handler).OnRead(string(message))
	}

}
