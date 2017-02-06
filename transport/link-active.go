package transport

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

// LinkActiver is an interface for connection struct.
// Realization of this interface depends of connection state:
// * connected
// * reconnecting
// * error state
type LinkActiver interface {
	Write(string) error
	Read() error
	GetLinkActiveID() string
	Disconnect() error
}

// LinkStub initialize LinkStater interface when
// connecting or reconnecting
type LinkStub struct {
	commandCh    *chan string
	LinkActiveID string
}

func (lAStub *LinkStub) Write(msg string) error {
	log.Warn("(lSub *LinkStub)Write")
	return fmt.Errorf("%s", ErrStubFunction)
}

func (lAStub *LinkStub) GetLinkActiveID() string {
	return lAStub.LinkActiveID
}

func (lAStub *LinkStub) Read() error {
	log.Warn("(lSub *LinkStub)Read")
	return fmt.Errorf("%s", ErrStubFunction)
}

func (lAStub *LinkStub) Disconnect() error {
	log.Info("(lSub *LinkStub) Disconnect")
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

func (lAWork *LinkActive) GetLinkActiveID() string {
	return lAWork.LinkActiveID
}

func (lAWork *LinkActive) Disconnect() error {
	log.Info("(lSub *LinkActive)Disconnect")
	return nil
}

func (lAWork *LinkActive) Write(msg string) error {

	log.Info("(lSub *LinkActive)Write")
	(*lAWork.conn).Write([]byte(msg + "\n"))
	return nil
}

//func (lAWork *LinkActive) Read(b []byte) error {
func (lAWork *LinkActive) Read() error {
	log.Info("(lSub *LinkActive)Read")

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
