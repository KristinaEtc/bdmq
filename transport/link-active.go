package transport

import (
	"net"

	"github.com/ventu-io/slf"
)

// LinkActive initialize LinkStater interface when
// connecting is ok
type LinkActive struct {
	conn         net.Conn
	linkDesc     *LinkDesc
	LinkActiveID string
	handler      Handler
	//parser       *parser.Parser
	commandCh      chan cmdActiveLink
	linkControl    *LinkControl
	log            slf.Logger
	FrameProcessor FrameProcessor
}

func (lA *LinkActive) Id() string {
	return lA.LinkActiveID
}

func (lA *LinkActive) Conn() net.Conn {
	return lA.conn
}

func (lA *LinkActive) Mode() int {
	return lA.linkControl.Mode()
}

func (lA *LinkActive) getHandler() Handler {
	return lA.handler
}

func (lA *LinkActive) Close() {
	lA.log.Info("Close()")
	lA.commandCh <- cmdActiveLink{
		cmd: quitLinkActive,
	}
}

func (lA *LinkActive) SendMessageActive(msg []byte) {
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
				lA.log.Debugf("got command=%s", command.cmd.String())
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
					lA.Write(command.msg)
				}
			}
		}
	}
}

func (lA *LinkActive) Write(msg []byte) error {

	lA.log.Warn("dProcessor Write")
	lA.conn.Write(msg)
	return nil
}

//func (lA *LinkActive) Read(b []byte) error {
func (lA *LinkActive) Read() {
	lA.log.Debug("Read")

	err := lA.FrameProcessor.Read()
	if err != nil {
		lA.linkControl.NotifyErrorRead(err)
		lA.linkControl.log.Warn("exiting")
	}
	//lA.Handler.OnRead(message)
	//msgStr := strings.TrimSpace(string(message))
	//lA.log.Debugf("Message Received: %s", msgStr)
	//	lA.Handler.OnRead(message)
}
