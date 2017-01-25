package transport

import (
	"net"
	"sync"
	"time"

	"github.com/KristinaEtc/bdmq/parser"
)

const network = "tcp"

const backOffLimit = time.Duration(time.Minute * 10)

// status enum
const (
	statusOnline       = iota
	statusOffline      = iota
	statusReconnecting = iota
)

// LinkDesc contains options for connecting to an address.
// It is store options for creating connection for ctrtain connection.
type LinkDesc struct {
	linkID  string
	address string
	mode    string //client/server
	handler string
}

// LinkActive provides a TCP connection with auto-reconnect capabilities.
type LinkActive struct {
	conn         *net.Conn
	lock         sync.RWMutex
	status       int32
	LinkConf     *LinkDesc
	LinkActiveID string
	handler      *Handler
	parser       *parser.Parser
}

// Node is a struct which combine connections with common Node's options.
type Node struct {
	NodeID      string
	LinkDesc    map[string]*LinkDesc
	LinkActives map[string]*LinkActive
}

//--------------------------------------

// NewNode creates an instance of Node struct.
func NewNode() (n *Node) {
	n = &Node{
		LinkDesc:    make(map[string]*LinkDesc),
		LinkActives: make(map[string]*LinkActive),
	}
	return
}

// InitLinkDesc parses config with nodes' options and creates nodes with such
// options.
// InitLinkDesc is a method of Noder interface from transport library.
func (n *Node) InitLinkDesc(lDescJSON []LinkDescFromJSON) error {

	log.Debug("InitLinkDesc")

	for _, l := range lDescJSON {

		lDesc := &LinkDesc{
			address: l.Address,
			linkID:  l.LinkID,
			mode:    l.Mode,
			handler: l.Handler,
		}

		n.LinkDesc[l.LinkID] = lDesc
	}

	return nil
}

//Run reads links and creates nodes with their options.
func (n *Node) Run() error {

	log.Debug("func Run")

	//	var err error

	if n.LinkDesc == nil && len(n.LinkDesc) == 0 {
		log.Debug(ErrEmptyLinkSlice.Error())
		return ErrEmptyLinkSlice
	}
	/*
		for _, linkConf := range n.LinkDesc {

			switch strings.ToLower(linkConf.mode) {
			//case "client":
			//	newConn, err = n.initClientLink(newLink, linkConf)
			case "server":
				go n.initServerLink(linkConf)
			default:
				log.Error(strings.ToLower(linkConf.mode))
				log.Warnf("%s (ID=%s)\n", ErrWrongNodeMode.Error(), linkConf.linkID)
			}
		}
	*/
	log.Debug("func Run closing")
	return nil
}

/*
func initHandler() {

}
*/

/*
func (n *Node) initLinkActive() {
	h := initHandler()
	for {
		msg := readConn.Read(b)

		errr = h.OnConn(msg) //проверяется пароль
		if err == nil {
			//добавить хендлер в activeLink
			// activeLink добавить список activeLinks в node
		}

	}
}
*/

/*
func (n *Node) initServerLink(linkConf *LinkDesc) {

	log.Debug("Init server node")

	//var conn net.Conn
	var err error
	var ln net.Listener
	for {
		// TODO: change this block when refactoring
		// listen on all interfaces
		ln, err = net.Listen(network, linkConf.address)
		if err != nil {
			log.WithField("ID=", linkConf.linkID).Errorf("Error listen: %s", err.Error())
			time.Sleep(time.Second * 1)
			continue
		}

		// accept connection on port
		conn, err = ln.Accept()
		if err != nil {
			log.Errorf("Error accept: %s", err.Error())
			time.Sleep(time.Second * 1)
			continue
		}

		/*
			//n.InitLinkActive()
			newActiveLink := LinkActive{
				conn: &conn,
			}

	}
}
*/

/*

	log.Debugf("Starting link ID=[%s]\n", linkConf.linkID)

	newLink := &TCPReconn{
		lock:         sync.RWMutex{},
		status:       0,
		LinkConf:     linkConf,
		LinkActiveID: linkConf.linkID + ":" + linkConf.address,
	}

	n.LinkActives = append(n.LinkActives, &TCPReconn{})


	if err != nil {
		log.Warnf("Could not connect: id=%s\n", linkConf.linkID)
		continue
	}
*/
