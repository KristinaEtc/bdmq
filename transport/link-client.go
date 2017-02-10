package transport

import (
	"net"
	"time"
)

type LinkControlClient struct {
	linkDesc  *LinkDesc
	node      *Node
	commandCh chan cmdContrlLink
	conn      net.Conn
}

func (lC *LinkControlClient) getId() string {
	return lC.linkDesc.linkID
}

func (lC *LinkControlClient) getChannel() chan cmdContrlLink {
	return lC.commandCh
}

func (lC *LinkControlClient) getLinkDesc() *LinkDesc {
	return lC.linkDesc
}

func (lC *LinkControlClient) getNode() *Node {
	return lC.node
}

func (lC *LinkControlClient) Close() {
	log.Info("(lC *LinkControlClient) Close()")
	lC.commandCh <- cmdContrlLink{
		cmd: quitControlLink,
	}
}

func (lC *LinkControlClient) NotifyError(err error) {
	lC.commandCh <- cmdContrlLink{
		cmd: errorControlLink,
		err: "Error" + err.Error(),
	}
}

/*
func (lC *LinkControlClient) NotifyErrorFromActive(err error) {
	lC.commandCh <- cmdContrlLink{
		cmd: errorFromActiveLink,
		err: "Error" + err.Error(),
	}
}
*/

func (lC *LinkControlClient) Dial() (net.Conn, error) {
	log.Debug("func Dial()")

	var err error
	var conn net.Conn

	var numOfRecon = 0
	var secToRecon = time.Duration(time.Second * 2)

	for {
		conn, err = net.Dial(network, lC.linkDesc.address)
		if err == nil {
			log.WithField("ID=", lC.linkDesc.linkID).Debugf("Established connection with: %s", conn.RemoteAddr().String())
			//lC.InitLinkActive(conn)
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
				if command.cmd == quitControlLink {
					log.WithField("ID=", lC.linkDesc).Info("Got quit commant. Closing link")
					//return nil, ErrQuitLinkRequested
					return nil, err
				}
				log.WithField("ID=", lC.linkDesc.linkID).Warnf("Got impermissible command %s. Ignored.", command)
			}
		}
	}
}

func (lC *LinkControlClient) WaitCommand() (isExiting bool) {

	select {
	case command := <-lC.commandCh:
		conn := lC.conn
		if command.cmd == quitControlLink {
			log.Debug("linkControlClient: quit received")
			if conn != nil {
				conn.Close()
			}
			return true
		}
		if command.cmd == errorControlLink {
			log.Errorf("Error: %s", command.err)
			/*if conn != nil {
				conn.Close()
			}*/
			return false
		}
	}
	log.Debug("linkControlClient: received something weird")
	return false
}

func (lC *LinkControlClient) InitLinkActive(conn net.Conn) {
	log.Debug("func InitLinkActive (client)")
	initLinkActive(lC, conn)

}
