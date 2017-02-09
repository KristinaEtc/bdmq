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
	commandCh chan cmdActiveLink
}

func (lA *LinkActive) Id() string {
	return lA.LinkActiveID
}

func (lA *LinkActive) Close() {
	log.Info("(func LinkActive Close()")
	lA.commandCh <- cmdActiveLink{
		cmd: quitLinkActive,
	}
}

func (lA *LinkActive) SendMessage(msg string) {
	lA.commandCh <- cmdActiveLink{
		cmd: sendMessageActive,
		msg: msg,
	}
}

func (lA *LinkActive) WaitCommand(conn net.Conn) (isExiting bool) {
	select {
	case command := <-lA.commandCh:
		{
			log.Debugf("got command=%s", command.cmd.String())
			if command.cmd == quitLinkActive {
				conn.Close()
				return true
			}
			if command.cmd == errorReadingActive {
				if lA.linkDesc.mode == "server" {
					conn.Close()
					return true
				}
			}
			if command.cmd == sendMessageActive {
				lA.handler.OnWrite(command.msg)
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
			lA.commandCh <- cmdActiveLink{
				cmd: errorReadingActive,
				err: err.Error(),
			}
			return
		}
		log.WithField("ID=", lA.LinkActiveID).Debugf("Message Received: %s", string(message))
		go lA.handler.OnRead(string(message))
	}
}
