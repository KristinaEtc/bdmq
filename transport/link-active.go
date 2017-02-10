package transport

import (
	"bufio"
	"net"
	"strings"
)

// LinkActive initialize LinkStater interface when
// connecting is ok
type LinkActive struct {
	conn         net.Conn
	linkDesc     *LinkDesc
	LinkActiveID string
	handler      Handler
	//parser       *parser.Parser
	commandCh   chan cmdActiveLink
	linkControl LinkControl
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

func (lA *LinkActive) WaitCommand(conn net.Conn) {
	for {
		select {
		case command := <-lA.commandCh:
			{
				log.Debugf("got command=%s", command.cmd.String())
				if command.cmd == quitLinkActive {
					conn.Close()
					return
				}
				/*if command.cmd == errorReadingActive {
					conn.Close()
					return
					//lA.linkControl.NotifyError(errors.New("Error reading active "))
				}*/
				if command.cmd == sendMessageActive {
					lA.handler.OnWrite([]byte(command.msg))
				}
			}
		}
	}
}

func (lA *LinkActive) Write(msg []byte) error {

	log.Info("(lSub *LinkActive)Write")
	lA.conn.Write(msg)
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
			lA.linkControl.NotifyErrorRead(err)
			log.Warn("exiting")
			return
		}
		msgStr := strings.TrimSpace(string(message))
		log.WithField("ID=", lA.LinkActiveID).Debugf("Message Received: %s", msgStr)
		lA.handler.OnRead(message)
	}
}
