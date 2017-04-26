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
	commandCh    chan cmdActiveLink
	linkControl  *LinkControl
	log          slf.Logger
}

// ID returns LinkActive ID
func (lA *LinkActive) ID() string {
	return lA.LinkActiveID
}

// Conn returns net.Conn of LinkActive
func (lA *LinkActive) Conn() net.Conn {
	return lA.conn
}

// Mode returns mode of LinkActive (server or client)
func (lA *LinkActive) Mode() int {
	return lA.linkControl.Mode()
}

// GetHandler returns handler of LinkActive
func (lA *LinkActive) GetHandler() Handler {
	return lA.handler
}

// Close sends close command to AcliveLink
func (lA *LinkActive) Close() {
	lA.log.Info("Close()")
	lA.commandCh <- cmdActiveLink{
		cmd: quitLinkActive,
	}
}

// SendMessageActive sends command to AcliveLink to send message msg
func (lA *LinkActive) SendMessageActive(msg []byte) {
	lA.commandCh <- cmdActiveLink{
		cmd: sendMessageActive,
		msg: msg,
	}
}

// RegisterTopic sends command to AcliveLink to create topic channel
func (lA *LinkActive) RegisterTopic(topicName string) {
	lA.commandCh <- cmdActiveLink{
		cmd: registerTopicActive,
		msg: []byte(topicName),
	}
}

// WaitCommand is a loop for LinkActive command processing
func (lA *LinkActive) WaitCommand(conn net.Conn) {
	for {
		select {
		case command := <-lA.commandCh:
			{
				lA.log.Debugf("Command=[%s/%v]", command.cmd.String(), command.cmd)
				if command.cmd == quitLinkActive {
					log.Debug("quitLinkActive")
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
				if command.cmd == registerTopicActive {

				}
			}
		}
	}
}

// Write gets msg and writes it to connection. It is an implementation of
// of method of interface LinkWriter.
func (lA *LinkActive) Write(msg []byte) (int, error) {

	//	lA.log.Debug("LinkActive Write() enter")
	n, err := lA.conn.Write(msg)
	return n, err
}

// Read calls method Read() from frameProcess interface.
// It is an implementation of method of interface LinkWriter.
func (lA *LinkActive) Read() {

	lA.log.Debug("LinkActive Read()")
	/*
		err := lA.FrameProcessor.Read()
		if err != nil {
			lA.linkControl.NotifyErrorRead(err)
			lA.linkControl.log.Warn("LinkActive.Read() exiting")
		}
	*/
	err := lA.GetHandler().OnRead(lA.conn)
	if err != nil {
		lA.linkControl.NotifyErrorRead(err)
		lA.log.Warn("LinkActive Read(): exiting")
	}

	//lA.Handler.OnRead(message)
	//msgStr := strings.TrimSpace(string(message))
	//lA.log.Debugf("Message Received: %s", msgStr)
	//	lA.Handler.OnRead(message)
}
