package transport

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/ventu-io/slf"
)

type LinkControl struct {
	linkDesc  *LinkDesc
	node      *Node
	commandCh chan cmdContrlLink
	isExiting bool
	log       slf.Logger
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
	lc.log.Debugf("sendCommand: %d %d", len(lc.commandCh), cap(lc.commandCh))
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

func (lc *LinkControl) Close() {
	lc.log.Info("Close()")
	lc.sendCommand(cmdContrlLink{
		cmd: quitControlLink,
	})
}

func (lc *LinkControl) NotifyErrorAccept(err error) {
	lc.sendCommand(cmdContrlLink{
		cmd: errorControlLinkAccept,
		err: "Error" + err.Error(),
	})
}

func (lc *LinkControl) NotifyErrorRead(err error) {
	//WEIRD HACK
	if lc.getLinkDesc().mode == "server" {
		return
	}
	lc.sendCommand(cmdContrlLink{
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

func (lc *LinkControl) Listen() (net.Listener, error) {
	lc.log.Debug("func Listen()")

	var err error
	var ln net.Listener

	var secToRecon = time.Duration(time.Second * 2)
	var numOfRecon = 0

	for {
		if lc.isExiting {
			return nil, fmt.Errorf("exiting")
		}
		ln, err = net.Listen(network, lc.linkDesc.address)
		if err == nil {
			lc.log.Debugf("Created listener with: %s", ln.Addr().String())
			return ln, nil
		}

		lc.log.Errorf("Error listen: %s. Reconnecting after %d milliseconds.", err.Error(), secToRecon/1000000.0)
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

		case command := <-lc.commandCh:
			{
				lc.log.Debugf("func Listen(): got command %s", command.cmd)
				if command.cmd == quitControlLink {
					lc.isExiting = true
					lc.log.Info("Got quit commant. Closing link")
					return nil, ErrQuitLinkRequested
				}
				lc.log.Warnf("Got impermissible command %s. Ignored.", command)
			}
		}
	}
}

func (lc *LinkControl) Accept(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			lc.log.Errorf("Error accept: %s", err.Error())
			lc.NotifyErrorAccept(err)
			return
		}
		lc.log.Debugf("Accept: new client %s", conn.RemoteAddr().String())

		go lc.initLinkActive(conn)
	}
}

func (lc *LinkControl) Dial() (net.Conn, error) {
	lc.log.Debug("func Dial()")

	var err error
	var conn net.Conn

	var numOfRecon = 0
	var secToRecon = time.Duration(time.Second * 2)

	for {
		if lc.isExiting {
			return nil, fmt.Errorf("exiting")
		}
		conn, err = net.Dial(network, lc.linkDesc.address)
		if err == nil {
			lc.log.Debugf("Established connection with: %s", conn.RemoteAddr().String())
			//lC.InitLinkActive(conn)
			return conn, nil
		}
		lc.log.Errorf("Error dial: %s. Reconnecting after %d milliseconds", err.Error(), secToRecon/1000000.0)
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

		case command := <-lc.commandCh:
			{
				log.Debugf("Dial: got command %+v", command)
				if command.cmd == quitControlLink {
					lc.isExiting = true
					lc.log.Info("Dial: Got quit commant. Closing link")
					//return nil, ErrQuitLinkRequested
					return nil, err
				}
				lc.log.Warnf("Got unknown command %+v. Ignored.", command)
			}
		}
	}
}

func (lc *LinkControl) WaitCommandServer(conn io.Closer) (isExiting bool) {
	for {
		select {
		case command := <-lc.commandCh:
			if command.cmd == quitControlLink {
				lc.log.Debugf("linkControl: quit received")
				lc.isExiting = true
				conn.Close()
				lc.WaitExit()

				return true
			}
			if command.cmd == errorControlLinkAccept {
				lc.log.Errorf("WaitCommandServer: error accept %s", command.err)
				conn.Close()
				return false
			}
			if command.cmd == errorControlLinkRead {
				lc.log.Errorf("WaitCommandServer: error read %s", command.err)
				//continue
			}
		}
	}

}

func (lc *LinkControl) WaitExit() {
	lc.log.Debugf("WaitExit")
	timeout := time.NewTimer(time.Duration(30) * time.Second)
	for {

		select {
		case _ = <-timeout.C:
			return
		case command := <-lc.commandCh:
			lc.log.Debugf("WaitExit: cmd %+v", command)
			return
		}
	}
}

func (lc *LinkControl) WaitCommandClient(conn io.Closer) (isExiting bool) {

	select {
	case command := <-lc.commandCh:
		if command.cmd == quitControlLink {
			lc.log.Debugf("WaitCommandClient: quit received")
			lc.isExiting = true
			conn.Close()
			lc.WaitExit()
			return true
		}
		if command.cmd == errorControlLinkRead {
			lc.log.Errorf("WaitCommandClient: error read %s", command.err)
			conn.Close()
			return false
		}
		lc.log.Warnf("WaitCommandClient: unknown command %+v", command)
		conn.Close()
		return false
	}
}

func (lc *LinkControl) initLinkActive(conn net.Conn) {
	id := lc.getId() + ":" + conn.RemoteAddr().String()
	lc.log.Debugf("InitLinkActive: %s", id)
	linkActive := LinkActive{
		conn:         conn,
		linkDesc:     lc.getLinkDesc(),
		LinkActiveID: id,
		commandCh:    make(chan cmdActiveLink),
		linkControl:  lc,
		log:          slf.WithContext("LinkActive").WithFields(slf.Fields{"ID": id}),
	}

	handlerName := lc.getLinkDesc().handler
	linkActive.log.Debugf("handler: %s", handlerName)

	hFactory, ok := handlers[handlerName]
	if !ok {
		linkActive.log.Errorf("initLinkActive: handler %s not found, closing connection", handlerName)
		conn.Close()
		lc.NotifyErrorRead(fmt.Errorf("initLinkActive: handler %s not found, closing connection", handlerName))
		return
	}

	node := lc.getNode()
	h := hFactory.InitHandler(&linkActive, node)
	linkActive.handler = h

	node.RegisterLinkActive(&linkActive)

	go linkActive.WaitCommand(conn)
	linkActive.handler.OnConnect()
	linkActive.Read()
	//linkActive.handler.OnDisconnect()
	node.UnregisterLinkActive(&linkActive)

	linkActive.log.Debug("initLinkActive exiting")
}

/*func (lc *LinkControl) InitLinkActive(conn net.Conn) {
	lc.log.Debug("InitLinkActive")
	initLinkActive(lC, conn)
}*/

func (lc *LinkControl) WorkClient() {

	for {
		lc.log.Debug("WorkClient: -------------------------------------")
		conn, err := lc.Dial()
		if err != nil {
			lc.log.Errorf("WorkClient: dial error: %s", err.Error())
			break
		}

		go lc.initLinkActive(conn)
		isExiting := lc.WaitCommandClient(conn)
		if isExiting {
			lc.log.Debug("WorkClient: isExiting client")
			break
		}
		lc.log.Debug("WorkClient: reconnect")

	}

	lc.log.Debug("WorkClient exit")
}

func (lc *LinkControl) WorkServer() {

	for {
		lc.log.Debug("WorkServer: -------------------------------------")
		ln, err := lc.Listen()
		if err != nil {
			lc.log.Errorf("WorkServer: listen error %s", err.Error())
			return
		}
		go lc.Accept(ln)

		isExiting := lc.WaitCommandServer(ln)
		if isExiting {
			lc.log.Debug("WorkClient: isExiting server")
			break
		}
	}

	lc.log.Debug("WorkServer: exit")
}
