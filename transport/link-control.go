package transport

import (
	"fmt"
	"io"
	"net"
	"runtime/debug"
	"time"

	"github.com/ventu-io/slf"
)

// LinkControl is a entity which represent configuration for new connections.
// One LinkControl can cause several LinkActives, which actually represent active connection.
type LinkControl struct {
	linkDesc  *LinkDesc
	node      *Node
	commandCh chan cmdContrlLink
	isExiting bool
	log       slf.Logger
	mode      int
	//conn      net.Conn
}

// Mode returnes a type of LinkControl: server or client.
func (lc *LinkControl) Mode() int {
	return lc.mode
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

func (lC *LinkControl) sendCommand(cmd cmdContrlLink) {
	lC.log.Debugf("sendCommand: %d %d", len(lC.commandCh), cap(lC.commandCh))
	/*for len(lc.commandCh) >= cap(lc.commandCh) {
		cmd := <-lc.commandCh
		log.Debugf("sendCommand: cleanup %+v", cmd)
	}*/
	lC.commandCh <- cmd
	/*if len(lc.commandCh) < cap(lc.commandCh) {
		lc.commandCh <- cmd
	} else {
		log.Warnf("sendCommand %s channel overflow", lc.getId())
	}*/

}

func (lC *LinkControl) Close() {
	lC.log.Debug("Close()")
	debug.PrintStack()
	lC.sendCommand(cmdContrlLink{
		cmd: quitControlLink,
	})
}

func (lC *LinkControl) NotifyErrorAccept(err error) {
	lC.sendCommand(cmdContrlLink{
		cmd: errorControlLinkAccept,
		err: err.Error(),
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
	lC.log.Debug("func Listen()")
	var err error
	var ln net.Listener

	var secToRecon = time.Duration(time.Second * 2)
	var numOfRecon = 0

	for {
		if lC.isExiting {
			return nil, fmt.Errorf("exiting")
		}
		ln, err = net.Listen(network, lC.linkDesc.address)
		if err == nil {
			lC.log.Debugf("Created listener with: %s", ln.Addr().String())
			return ln, nil
		}

		lC.log.Errorf("Error listen: %s. Reconnecting after %d milliseconds.", err.Error(), secToRecon/1000000.0)
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
				lC.log.Debugf("func Listen(): got command %s", command.cmd)
				if command.cmd == quitControlLink {
					lC.isExiting = true
					lC.log.Debug("Got quit commant. Closing link")
					return nil, ErrQuitLinkRequested
				}
				lC.log.Warnf("Got impermissible command %s. Ignored.", command)
			}
		}
	}
}

func (lC *LinkControl) Accept(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			lC.log.Errorf("Error accept: %s", err.Error())
			lC.NotifyErrorAccept(err)
			return
		}
		lC.log.Debugf("Accept: new client %s", conn.RemoteAddr().String())

		go lC.initLinkActive(conn)
	}
}

func (lC *LinkControl) Dial() (net.Conn, error) {
	lC.log.Debug("func Dial()")

	var err error
	var conn net.Conn

	var numOfRecon = 0
	var secToRecon = time.Duration(time.Second * 2)

	for {
		if lC.isExiting {
			return nil, fmt.Errorf("exiting")
		}
		conn, err = net.Dial(network, lC.linkDesc.address)
		if err == nil {
			lC.log.Debugf("Established connection with: %s", conn.RemoteAddr().String())
			//lC.InitLinkActive(conn)
			return conn, nil
		}
		lC.log.Errorf("Error dial: %s. Reconnecting after %d milliseconds", err.Error(), secToRecon/1000000.0)
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
				log.Debugf("Dial: got command %+v", command)
				if command.cmd == quitControlLink {
					lC.isExiting = true
					lC.log.Debug("Dial: Got quit commant. Closing link")
					//return nil, ErrQuitLinkRequested
					return nil, err
				}
				lC.log.Warnf("Got unknown command %+v. Ignored.", command)
			}
		}
	}
}

func (lC *LinkControl) WaitCommandServer(conn io.Closer) (isExiting bool) {
	for {
		select {
		case command := <-lC.commandCh:
			if command.cmd == quitControlLink {
				lC.log.Debugf("WaitCommandServer: quit received")
				lC.isExiting = true
				conn.Close()
				lC.WaitExit()

				return true
			}
			if command.cmd == errorControlLinkAccept {
				lC.log.Errorf("WaitCommandServer: error accept %s", command.err)
				conn.Close()
				return false
			}
			if command.cmd == errorControlLinkRead {
				lC.log.Errorf("WaitCommandServer: error read %s", command.err)
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
			lC.log.Debugf("WaitExit: cmd %+v", command)
			return
		}
	}
}

func (lC *LinkControl) WaitCommandClient(conn io.Closer) (isExiting bool) {

	select {
	case command := <-lC.commandCh:
		if command.cmd == quitControlLink {
			lC.log.Debugf("linkControl: quit received %s", lC.getId())
			lC.isExiting = true
			conn.Close()
			lC.WaitExit()
			return true
		}
		if command.cmd == errorControlLinkRead {
			lC.log.Errorf("Error: %s", command.err)
			conn.Close()
			return false
		}
		lC.log.Warnf("linkControl: received something weird %d", command.cmd)
		conn.Close()
		return false
	}
}

func (lC *LinkControl) initLinkActive(conn net.Conn) {

	id := lC.getId() //+ ":" + conn.LocalAddr().String() + "-" + conn.RemoteAddr().String()
	lC.log.Debugf("InitLinkActive: %s", id)

	linkActive := LinkActive{
		conn:         conn,
		linkDesc:     lC.getLinkDesc(),
		LinkActiveID: id,
		commandCh:    make(chan cmdActiveLink),
		linkControl:  lC,
		log:          slf.WithContext("LinkActive").WithFields(slf.Fields{"ID": id}),
	}

	linkActive.log.Debugf("frame processer: %s", lC.getLinkDesc().frameProcessor)

	frameProcessorFactory, ok := frameProcessors[lC.linkDesc.frameProcessor]
	if !ok {
		linkActive.log.Warn("initLinkActive: frame processor not found, will be used default")
		linkActive.FrameProcessor = dFrameProcessorFactory.InitFrameProcessor(linkActive, conn, conn, linkActive.log)
	} else {
		linkActive.FrameProcessor = frameProcessorFactory.InitFrameProcessor(linkActive, conn, conn, linkActive.log)
	}

	handlerName := lC.getLinkDesc().handler
	linkActive.log.Debugf("handler: %s", handlerName)

	hFactory, ok := handlers[handlerName]
	if !ok {
		linkActive.log.Errorf("initLinkActive: handler %s not found, closing connection", handlerName)
		conn.Close()
		lC.NotifyErrorRead(fmt.Errorf("initLinkActive: handler %s not found, closing connection", handlerName))
		return
	}

	node := lC.getNode()
	h := hFactory.InitHandler(&linkActive, node)
	linkActive.handler = h

	node.RegisterLinkActive(&linkActive)

	go linkActive.WaitCommand(conn)
	linkActive.Read()
	//linkActive.handler.OnDisconnect()

	node.UnregisterLinkActive(&linkActive)
	linkActive.log.Debug("initLinkActive exiting")
}

/*
func (lC *LinkControl) InitLinkActive(conn net.Conn) {
	log.Debug("func InitLinkActive")
	initLinkActive(lC, conn)
}
*/

func (lC *LinkControl) WorkClient() {

	for {
		lC.log.Debug("WorkClient: -------------------------------------")
		conn, err := lC.Dial()
		if err != nil {
			lC.log.Errorf("WorkClient: dial error: %s", err.Error())
			break
		}

		go lC.initLinkActive(conn)
		isExiting := lC.WaitCommandClient(conn)
		if isExiting {
			lC.log.Debug("WorkClient: isExiting client")
			break
		}
		lC.log.Debug("WorkClient: reconnect")
	}

	lC.log.Debug("WorkClient exit")
}

func (lC *LinkControl) WorkServer() {

	for {
		lC.log.Debug("WorkServer: -------------------------------------")
		ln, err := lC.Listen()
		if err != nil {
			lC.log.Errorf("WorkServer: listen error %s", err.Error())
			return
		}
		go lC.Accept(ln)

		isExiting := lC.WaitCommandServer(ln)
		if isExiting {
			lC.log.Debug("WorkClient: isExiting server")
			break
		}
	}

	lC.log.Debug("WorkServer: exit")
}
