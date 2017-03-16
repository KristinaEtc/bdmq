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
func (lC *LinkControl) Mode() int {
	return lC.mode
}

func (lC *LinkControl) getID() string {
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

// Close sends command to close ControlLink
func (lC *LinkControl) Close() {
	lC.log.Debug("Close()")
	debug.PrintStack()
	lC.sendCommand(cmdContrlLink{
		cmd: quitControlLink,
	})
}

// NotifyErrorAccept sends command to notify an error in accept
// to he handles it by LinkControl
func (lC *LinkControl) NotifyErrorAccept(err error) {
	lC.sendCommand(cmdContrlLink{
		cmd: errorControlLinkAccept,
		err: err.Error(),
	})
}

// NotifyErrorRead sends command to notify an error in read
// to he handles it by LinkControl
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

// Listen runs listen loop and waiting new connections (server) of trying to make a new connection (client)/
// If new connection was esteblished, executed
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
					return nil, errQuitLinkRequested
				}
				lC.log.Warnf("Got impermissible command %s. Ignored.", command)
			}
		}
	}
}

// Accept runs an the Accemt method of the net.Listener interface; it waits for the next call and
// if connection established successfuly, calls initLinkActive to create new LinkActive.
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

// Dial is a method which calls net.Conn Dial() in a loop till connection established successfully
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

// WaitCommandServer processes commands of LinkControl with mode "server"
func (lC *LinkControl) WaitCommandServer(conn io.Closer) (isExiting bool) {
	for {
		select {
		case command := <-lC.commandCh:
			if command.cmd == quitControlLink {
				lC.log.Debugf("WaitCommandServer: quit received")
				lC.isExiting = true
				conn.Close()
				lC.waitExit()

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

func (lC *LinkControl) waitExit() {
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

// WaitCommandClient processes commands of LinkControl with mode "client"
func (lC *LinkControl) WaitCommandClient(conn io.Closer) (isExiting bool) {

	select {
	case command := <-lC.commandCh:
		if command.cmd == quitControlLink {
			lC.log.Debugf("linkControl: quit received %s", lC.getID())
			lC.isExiting = true
			conn.Close()
			lC.waitExit()
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

	id := lC.getID() //+ ":" + conn.LocalAddr().String() + "-" + conn.RemoteAddr().String()
	lC.log.Debugf("InitLinkActive: %s", id)

	linkActive := LinkActive{
		conn:         conn,
		linkDesc:     lC.getLinkDesc(),
		LinkActiveID: id,
		commandCh:    make(chan cmdActiveLink),
		linkControl:  lC,
		log:          slf.WithContext("LinkActive").WithFields(slf.Fields{"ID": id}),
	}

	frameProcessorFactory, ok := frameProcessors[lC.linkDesc.frameProcessor]
	if !ok {
		linkActive.log.Warnf("initLinkActive: frame processor %s not found, will be used default", frameProcessorFactory)
		linkActive.FrameProcessor = dFrameProcessorFactory.initFrameProcessor(linkActive, conn, conn)
	} else {
		linkActive.FrameProcessor = frameProcessorFactory.InitFrameProcessor(linkActive, conn, conn)
	}

	handlerName := lC.getLinkDesc().handler

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

// WorkClient runs Dial() to current LinkControl, then, if connection was established, creates new LinkActive
// and calls WaitCommandServer with loop which process commands for this
// LinkActive.
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

// WorkServer runs Listen() to current LinkControl, then, if connection was established, Accept(),
// where a new LinkActive will be created. Then runs WaitCommandServer with loop which process commands for this
// LinkActive.
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
