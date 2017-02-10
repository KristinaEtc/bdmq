package transport

import (
	"fmt"
	"io"
	"net"
	"time"
)

type LinkControl struct {
	linkDesc  *LinkDesc
	node      *Node
	commandCh chan cmdContrlLink
	//conn      net.Conn
}

func (lC *LinkControl) getId() string {
	return lC.linkDesc.linkID
}

/*func (lC *LinkControl) getChannel() chan cmdContrlLink {
	return lC.commandCh
}*/

func (lC *LinkControl) getLinkDesc() *LinkDesc {
	return lC.linkDesc
}

func (lC *LinkControl) getNode() *Node {
	return lC.node
}

func (lc *LinkControl) sendCommand(cmd cmdContrlLink) {
	log.Debugf("sendCommand: %d %d", len(lc.commandCh), cap(lc.commandCh))
	/*for len(lc.commandCh) >= cap(lc.commandCh) {
		cmd := <-lc.commandCh
		log.Debugf("sendCommand: cleanup %+v", cmd)
	}*/
	lc.commandCh <- cmd
	/*if len(lc.commandCh) < cap(lc.commandCh) {
		lc.commandCh <- cmd
	} else {
		log.Warnf("sendCommand %s channel overflow", lc.getId())
	}*/

}

func (lC *LinkControl) Close() {
	log.Info("(lC *LinkControl) Close()")
	lC.sendCommand(cmdContrlLink{
		cmd: quitControlLink,
	})
}

func (lC *LinkControl) NotifyErrorAccept(err error) {
	lC.sendCommand(cmdContrlLink{
		cmd: errorControlLinkAccept,
		err: "Error" + err.Error(),
	})
}

func (lC *LinkControl) NotifyErrorRead(err error) {
	//WEIRD HACK
	if lC.getLinkDesc().mode == "server" {
		return
	}
	lC.sendCommand(cmdContrlLink{
		cmd: errorControlLinkRead,
		err: "Error" + err.Error(),
	})
}

/*
func (lC *LinkControlClient) NotifyErrorFromActive(err error) {
	lC.commandCh <- cmdContrlLink{
		cmd: errorFromActiveLink,
		err: "Error" + err.Error(),
	}
}
*/

func (lC *LinkControl) Listen() (net.Listener, error) {
	log.Debug("func Listen()")

	var err error
	var ln net.Listener

	var secToRecon = time.Duration(time.Second * 2)
	var numOfRecon = 0

	for {
		ln, err = net.Listen(network, lC.linkDesc.address)
		if err == nil {
			log.WithField("ID=", lC.linkDesc.linkID).Debugf("Created listener with: %s", ln.Addr().String())
			return ln, nil
		}

		log.WithField("ID=", lC.linkDesc.linkID).Errorf("Error listen: %s. Reconnecting after %d milliseconds.", err.Error(), secToRecon/1000000.0)
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
				log.Debugf("func Listen(): got command %s", command.cmd)
				if command.cmd == quitControlLink {
					log.WithField("ID=", lC.linkDesc).Info("Got quit commant. Closing link")
					return nil, ErrQuitLinkRequested
				}
				log.WithField("ID=", lC.linkDesc.linkID).Warnf("Got impermissible command %s. Ignored.", command)
			}
		}
	}
}

func (lC *LinkControl) Accept(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.WithField("ID=", lC.linkDesc.linkID).Errorf("Error accept: %s", err.Error())
			lC.NotifyErrorAccept(err)
			return
		}
		log.WithField("ID=", lC.linkDesc.linkID).Debug("New client")

		go lC.InitLinkActive(conn)
	}
}

func (lC *LinkControl) Dial() (net.Conn, error) {
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

func (lC *LinkControl) WaitCommandServer(conn io.Closer) (isExiting bool) {
	for {
		select {
		case command := <-lC.commandCh:
			if command.cmd == quitControlLink {
				log.Debug("linkControl: quit received")
				conn.Close()
				lC.WaitExit()

				return true
			}
			if command.cmd == errorControlLinkAccept {
				log.Error(command.err)
				conn.Close()
				return false
			}
			if command.cmd == errorControlLinkRead {
				log.Error(command.err)
				//continue
			}
		}
	}

}

func (lC *LinkControl) WaitExit() {
	log.Debugf("WaitExit")
	timeout := time.NewTimer(time.Duration(30) * time.Second)
	for {

		select {
		case _ = <-timeout.C:
			return
		case command := <-lC.commandCh:
			log.Debugf("WaitExit: cmd %+v", command)
			return
		}
	}
}

func (lC *LinkControl) WaitCommandClient(conn io.Closer) (isExiting bool) {

	select {
	case command := <-lC.commandCh:
		if command.cmd == quitControlLink {
			log.Debugf("linkControl: quit received %s", lC.getId())
			conn.Close()
			lC.WaitExit()
			return true
		}
		if command.cmd == errorControlLinkRead {
			log.Errorf("Error: %s", command.err)
			conn.Close()
			return false
		}
		log.Warnf("linkControl: received something weird %d", command.cmd)
		conn.Close()
		return false
	}
}

func initLinkActive(lCntl *LinkControl, conn net.Conn) {

	linkActive := LinkActive{
		conn:         conn,
		linkDesc:     lCntl.getLinkDesc(),
		LinkActiveID: lCntl.getId() + ":" + conn.RemoteAddr().String(),
		commandCh:    make(chan cmdActiveLink),
		linkControl:  *lCntl,
	}

	log.Debug(lCntl.getLinkDesc().handler)

	if _, ok := handlers[lCntl.getLinkDesc().handler]; !ok {
		log.Error("No handler! Closing linkControl")
		conn.Close()
		lCntl.NotifyErrorRead(fmt.Errorf("Error: " + "No handler! Closing linkControl"))
	}

	node := lCntl.getNode()
	h := handlers[lCntl.getLinkDesc().handler].InitHandler(&linkActive, node)
	linkActive.handler = h

	node.RegisterLinkActive(&linkActive)

	go linkActive.WaitCommand(conn)
	linkActive.handler.OnConnect()
	linkActive.Read()
	log.Debug("initLinkActive exiting")
	//linkActive.handler.OnDisconnect()

	node.UnregisterLinkActive(&linkActive)
}

func (lC *LinkControl) InitLinkActive(conn net.Conn) {
	log.Debug("func InitLinkActive")
	initLinkActive(lC, conn)
}

func (lC *LinkControl) WorkClient() {

	for {
		conn, err := lC.Dial()
		if err != nil {
			log.Errorf("Dial error: %s", err.Error())
			break
		}

		go lC.InitLinkActive(conn)
		isExiting := lC.WaitCommandClient(conn)
		if isExiting {
			log.Debug("isExiting client")
			break
		}
		log.Debug("reconnect")

	}

	log.Debugf("WorkClient exit %s", lC.getId())
}

func (lC *LinkControl) WorkServer() {

	for {
		ln, err := lC.Listen()
		if err != nil {
			log.Errorf("Listen error: %s", err.Error())
			return
		}
		go lC.Accept(ln)

		isExiting := lC.WaitCommandServer(ln)
		if isExiting {
			break
		}
	}

	log.Debugf("WorkServer exit %s", lC.getId())
}
