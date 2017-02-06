package transport

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

// LinkActive is an interface for connection struct.
// Realization of this interface depends of connection state:
// * connected
// * reconnecting
// * error state
type LinkerActive interface {
	Write(string) error
	Read() error
	GetLinkActiveID() string
	Disconnect() error
}

// LinkActiveStub initialize LinkActive interface when
// connecting or reconnecting
type LinkActiveStub struct {
	commandCh    *chan string
	LinkActiveID string
}

func (lAStub *LinkActiveStub) Write(msg string) error {
	log.Warn("(lSub *LinkActiveStub)Write")
	return fmt.Errorf("%s", ErrStubFunction)
}

func (lAStub *LinkActiveStub) GetLinkActiveID() string {
	return lAStub.LinkActiveID
}

func (lAStub *LinkActiveStub) Read() error {
	log.Warn("(lSub *LinkActiveStub)Read")
	return fmt.Errorf("%s", ErrStubFunction)
}

func (lAStub *LinkActiveStub) Disconnect() error {
	log.Info("(lSub *LinkActiveStub) Disconnect")
	return fmt.Errorf("%s", ErrStubFunction)
}

// LinkActiveWork initialize LinkActive interface when
// connecting is ok
type LinkActiveWork struct {
	conn         *net.Conn
	lock         sync.RWMutex
	status       int32
	LinkConf     *LinkDesc
	LinkActiveID string
	handler      *Handler
	//parser       *parser.Parser
	commandCh *chan string
}

func (lAWork *LinkActiveWork) GetLinkActiveID() string {
	return lAWork.LinkActiveID
}

func (lAWork *LinkActiveWork) Disconnect() error {
	log.Info("(lSub *LinkActiveStub)Disconnect")
	return nil
}

func (lAWork *LinkActiveWork) Write(msg string) error {

	log.Info("(lSub *LinkActiveStub)Write")
	(*lAWork.conn).Write([]byte(msg + "\n"))
	return nil
}

//func (lAWork *LinkActiveWork) Read(b []byte) error {
func (lAWork *LinkActiveWork) Read() error {
	log.Info("(lSub *LinkActiveWork)Read")

	buf := make([]byte, SizeOfBuf)

	for {
		message, err := bufio.NewReader(*lAWork.conn).Read(buf)
		if err != nil {
			log.WithField("ID=", lAWork.LinkActiveID).Errorf("Error read: %s", err.Error())
			log.WithField("ID=", lAWork.LinkActiveID).Errorf("ADD GRASEFUL RECONNECT")
			os.Exit(1)
		}
		log.WithField("ID=", lAWork.LinkActiveID).Debugf("Message Received: %s", string(message))
		go (*lAWork.handler).OnRead(string(message))
	}

}
