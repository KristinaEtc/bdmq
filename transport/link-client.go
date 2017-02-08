package transport

import (
	"net"
	"time"
)

type LinkControlClient struct {
	linkDesc  *LinkDesc
	node      *Node
	commandCh chan string
}

func (lC *LinkControlClient) Close() {
	log.Info("(lC *LinkControlClient) Close()")
	lC.commandCh <- "quit"
}

func (lC *LinkControlClient) Id() string {
	return lC.linkDesc.linkID
}

func (lC *LinkControlClient) NotifyError(err error) {
	lC.commandCh <- "Error" + err.Error()

}

func (lC *LinkControlClient) Dial() (net.Conn, error) {
	log.Debug("func Dial()")

	var err error
	var conn net.Conn

	var numOfRecon = 0
	var secToRecon = time.Duration(time.Second * 2)

	for {
		conn, err = net.Dial(lC.linkDesc.address, network)
		if err == nil {
			log.WithField("ID=", lC.linkDesc.linkID).Debugf("Established connection with: %s", conn.RemoteAddr().String())
			return conn, nil
		}
		log.WithField("ID=", lC.linkDesc.linkID).Errorf("Error dial: %s. Reconnecting after %d milliseconds", err.Error(), secToRecon/1000000.0)
		ticker := time.NewTicker(secToRecon)

		select {
		case _ = <-ticker.C:
			{
				if secToRecon < backOffLimit {
					randomAdd := secToRecon / 100 * (20 + time.Duration(r1.Int31n(10)))
					secToRecon = secToRecon*2 + time.Duration(randomAdd)
				}
				numOfRecon++
				ticker = time.NewTicker(secToRecon)
				continue
			}

		case command := <-lC.commandCh:
			{
				log.Debugf("func Dial(): got command %s", command)
				if command == "quit" {
					log.WithField("ID=", lC.linkDesc).Info("Got quit commant. Closing link")
					return nil, ErrQuitLinkRequested
				}
				log.WithField("ID=", lC.linkDesc.linkID).Warnf("Got impermissible command %s. Ignored.", command)
			}
		}
	}
}

func (lC *LinkControlClient) WaitCommand(conn net.Conn) (isExiting bool) {
	select {
	case command := <-lC.commandCh:
		if command == "quit" {
			log.Debug("linkControlClient: quit recieved")
			conn.Close()
			return true
		}
		return false
	}
}
