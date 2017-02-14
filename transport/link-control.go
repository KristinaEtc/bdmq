package transport

import (
	"fmt"
	"io"
	"net"
	"time"

	"runtime/debug"

	"github.com/ventu-io/slf"
)

type LinkControl struct {
	linkDesc  *LinkDesc
	mode      int
	node      *Node
	commandCh chan cmdContrlLink
	isExiting bool
	log       slf.Logger
	//conn      net.Conn
}

func (lc *LinkControl) getId() string {
	return lc.linkDesc.linkID
}

func (lc *LinkControl) Mode() int {
	return lc.mode
}

/*func (lC *LinkControl) getChannel() chan cmdContrlLink {
	return lC.commandCh
}*/

func (lc *LinkControl) getLinkDesc() *LinkDesc {
	return lc.linkDesc
}

func (lc *LinkControl) getNode() *Node {
	return lc.node
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
	lc.log.Debug("Close()")
	debug.PrintStack()
	lc.sendCommand(cmdContrlLink{
		cmd: quitControlLink,
	})
}

func (lc *LinkControl) NotifyErrorAccept(err error) {
	lc.sendCommand(cmdContrlLink{
		cmd: errorControlLinkAccept,
		err: err.Error(),
	})
}

func (lc *LinkControl) NotifyErrorRead(err error) {
	//WEIRD HACK
	if lc.getLinkDesc().mode == "server" {
		return
	}
	lc.sendCommand(cmdContrlLink{
		cmd: errorControlLinkRead,
		err: err.Error(),
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
					lc.log.Debug("Got quit commant. Closing link")
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

		go lC.InitActiveLink(conn)
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
					lc.log.Debug("Dial: Got quit commant. Closing link")
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
				lc.log.Debugf("WaitCommandServer: quit received")
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

func (lc *LinkControl) initActiveLink(conn net.Conn) {
	id := lc.getId() + ":" + conn.LocalAddr().String() + "-" + conn.RemoteAddr().String()
	lc.log.Debugf("InitLinkActive: %s", id)
	linkActive := ActiveLink{
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

	node := lCntl.getNode()
	h := handlers[lCntl.getLinkDesc().handler].InitHandler(&activeLink, node)
	activeLink.handler = h


	node.RegisterActiveLink(&activeLink)

	go 