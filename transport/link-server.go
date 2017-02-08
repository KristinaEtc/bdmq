package transport

import (
	"net"
	"time"

	//"github.com/KristinaEtc/bdmq1/transport"
)

type LinkControlServer struct {
	linkDesc  *LinkDesc
	node      *Node
	commandCh chan string
}

func (lS *LinkControlServer) Close() {
	log.Info("(lS *LinkControlServer) Close()")
	lS.commandCh <- "quit"
}

func (lS *LinkControlServer) Id() string {
	return (lS.linkDesc).linkID
}

func (lS *LinkControlServer) NotifyError(err error) {
	lS.commandCh <- "Error" + err.Error()

}

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
				log.Debugf("func Listen(): got command %s", command)
				if command == "quit" {
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
		if command == "quit" {
			log.Debug("linkControlServer: quit recieved")
			listener.Close()
			return true
		}
		return false
	}
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
	log.Debug("func InitLinkActive")

	linkActive := LinkActive{
		conn:         conn,
		linkDesc:     lS.linkDesc,
		LinkActiveID: lS.linkDesc.linkID + ":" + conn.RemoteAddr().String(),
		commandCh:    make(chan commandActiveLink),
	}

	log.Debug(lS.linkDesc.handler)

	if _, ok := handlers[lS.linkDesc.handler]; !ok {
		log.Error("No handler! Closing linkControl")
		lS.commandCh <- "quit"
	}
	h := handlers[lS.linkDesc.handler].InitHandler(&linkActive, lS.node)
	linkActive.handler = h

	lS.node.RegisterLinkActive(&linkActive)

	for {
		linkActive.handler.OnConnect()
		//linkActive.handler.Read()
		go linkActive.Read()
		isExiting := linkActive.WaitCommand(conn)
		if isExiting {
			break
		}
	}

	lS.node.UnregisterLinkActive(&linkActive)
}
