package transport

import (
	"net"
	"time"
)

type LinkControlServer struct {
	linkDesc  *LinkDesc
	node      *Node
	commandCh chan cmdContrlLink
}

func (lS *LinkControlServer) getId() string {
	return (lS.linkDesc).linkID
}

func (lS *LinkControlServer) getChannel() chan cmdContrlLink {
	return lS.commandCh
}

func (lS *LinkControlServer) getLinkDesc() *LinkDesc {
	return lS.linkDesc
}

func (lS *LinkControlServer) getNode() *Node {
	return lS.node
}

func (lS *LinkControlServer) Close() {
	log.Info("(lS *LinkControlServer) Close()")
	lS.commandCh <- cmdContrlLink{
		cmd: quitControlLink,
	}
}

func (lS *LinkControlServer) NotifyError(err error) {
	lS.commandCh <- cmdContrlLink{
		cmd: errorControlLink,
		err: "Error" + err.Error(),
	}
}

/*
func (lS *LinkControlServer) NotifyErrorFromActive(err error) {
	/*lS.commandCh <- cmdContrlLink{
		cmd: errorControlLink,
		err: "Error" + err.Error(),
	}
	log.Error(err.Error())
}
*/

func (lS *LinkControlServer) Listen() (net.Listener, error) {
	log.Debug("func Listen()")

	var err error
	var ln net.Listener

	var secToRecon = time.Duration(time.Second * 2)
	var numOfRecon = 0

	for {
		ln, err = net.Listen(network, lS.linkDesc.address)
		if err == nil {
			log.WithField("ID=", lS.linkDesc.linkID).Debugf("Created listener with: %s", ln.Addr().String())
			return ln, nil
		}

		log.WithField("ID=", lS.linkDesc.linkID).Errorf("Error listen: %s. Reconnecting after %d milliseconds.", err.Error(), secToRecon/1000000.0)
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

		case command := <-lS.commandCh:
			{
				log.Debugf("func Listen(): got command %s", command.cmd)
				if command.cmd == quitControlLink {
					log.WithField("ID=", lS.linkDesc).Info("Got quit commant. Closing link")
					return nil, ErrQuitLinkRequested
				}
				log.WithField("ID=", lS.linkDesc.linkID).Warnf("Got impermissible command %s. Ignored.", command)
			}
		}
	}
}

func (lS *LinkControlServer) WaitCommand(listener net.Listener) (isExiting bool) {
	select {
	case command := <-lS.commandCh:
		if command.cmd == quitControlLink {
			log.Debug("linkControlServer: quit recieved")
			listener.Close()
			return true
		}
		if command.cmd == errorControlLink {
			log.Error(command.err)
			listener.Close()
			return false
		}
	}
	return false
}

func (lS *LinkControlServer) Accept(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.WithField("ID=", lS.linkDesc.linkID).Errorf("Error accept: %s", err.Error())
			lS.NotifyError(err)
			return
		}
		log.WithField("ID=", lS.linkDesc.linkID).Debug("New client")

		lS.InitLinkActive(conn)
		conn.Close()
	}
}

func (lS *LinkControlServer) InitLinkActive(conn net.Conn) {
	log.Debug("func InitLinkActive (server)")
	initLinkActive(lS, conn)
}

func (lS *LinkControlServer) Work() {

	for {
		ln, err := lS.Listen()
		if err != nil {
			log.Errorf("Listen error: %s", err.Error())
			break
		}
		go lS.Accept(ln)
		isExiting := lS.WaitCommand(ln)
		if isExiting {
			break
		}
	}

}
