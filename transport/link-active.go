package transport

import (
	"bufio"
	"net"
)

// LinkActive initialize LinkStater interface when
// connecting is ok
type LinkActive struct {
	conn         net.Conn
	linkDesc     *LinkDesc
	LinkActiveID string
	handler      Handler
	//parser       *parser.Parser
	commandCh chan commandActiveLink
}

func (lA *LinkActive) Id() string {
	return lA.LinkActiveID
}

func (lA *LinkActive) Close() {
	log.Info("(func LinkActive Close()")
	lA.commandCh <- quitLinkActive
}

func (lA *LinkActive) WaitCommand(conn net.Conn) (isExiting bool) {
	select {
	case command := <-lA.commandCh:
		{
			if command == quitLinkActive {
				log.Debug("linkActive: quit recieved")
				conn.Close()
				return true
			}
			if command == errorReading {
				log.Debug("linkActive: errorReading recieved")
				if lA.linkDesc.mode == "server" {
					conn.Close()
					return true
				}
			}
		}
	}
	return false
}

func (lA *LinkActive) Write(msg string) error {

	log.Info("(lSub *LinkActive)Write")
	lA.conn.Write([]byte(msg + "\n"))
	return nil
}

//func (lA *LinkActive) Read(b []byte) error {
func (lA *LinkActive) Read() {
	log.Info("(lSub *LinkActive)Read")

	//	buf := make([]byte, lA.linkDesc.bufSize)

	for {
		message, err := bufio.NewReader(lA.conn).ReadBytes('\n')
		if err != nil {
			log.WithField("ID=", lA.LinkActiveID).Errorf("Error read: %s", err.Error())
			lA.commandCh <- errorReading
			return
		}
		log.WithField("ID=", lA.LinkActiveID).Debugf("Message Received: %s", string(message))
		go lA.handler.OnRead(string(message))
	}
}
