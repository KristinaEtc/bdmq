package tcprec

import (
	"encoding/json"
	"io"
	"math"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ventu-io/slf"
)

var log = slf.WithContext("tcprec.go")

const network = "tcp"

const backOffLimit = time.Duration(time.Minute * 10)

// status enum
const (
	statusOnline       = iota
	statusOffline      = iota
	statusReconnecting = iota
)

func (n *Node) parseConf(data []byte) error {

	err := json.Unmarshal(data, n.LinkDesc)
	if err != nil {
		return err
	}

	return nil
}

// LinkDesc contains options for connecting to an address.
// It is store options for creating connection for ctrtain connection.
type LinkDesc struct {
	linkID string
	// Address is an address to use when dialing an
	// address
	address  string
	mode     string
	Internal int
	lock     sync.RWMutex
}

// TCPReconn provides a TCP connection with auto-reconnect capabilities.
type TCPReconn struct {
	TCPConn      *net.TCPConn
	lock         sync.RWMutex
	status       int32
	LinkConf     *LinkDesc
	LinkActiveID string
}

// Node is a struct which conbine connections with common Node's options.
// Noder realizes Noder interface from transport library.
type Node struct {
	NodeID      string
	LinkDesc    []*LinkDesc
	LinkActives []*TCPReconn
}

//--------------------------------------

// NewNode creates an instance of Node struct.
func NewNode() (n *Node) {
	n = &Node{}
	return
}

// InitLinkDesc parses config with nodes' options and creates nodes with such
// options.
// InitLinkDesc is a method of Noder interface from transport library.
func (n *Node) InitLinkDesc(confJSON []byte) error {

	log.Debug("InitLinkDesc")

	err := n.parseConf(confJSON)
	if err != nil {
		log.Debugf("Initialize node from code config - err: %d", err.Error())
		return err
	}

	return nil
}

// AddLinkDesc adds link with specified options.
// AddLinkDesc is a method of Noder interface from transport library.
func (n *Node) AddLinkDesc(confJSON []byte) error {

	log.Debug("AddLinkDesc")

	// TODO: PARSECONF PARSED TO ARRAYS OF STRUCT, NOT
	// ONE STRUCT
	err := n.parseConf(confJSON)
	if err != nil {
		log.Debugf("Initialize node from code config - err: %d", err.Error())
		return err
	}

	return nil
}

//Run reads links and creates nodes with their options.
func (n *Node) Run() error {

	log.Debug("func Run")

	var err error
	var newConn net.Conn

	if n.LinkDesc == nil && len(n.LinkDesc) == 0 {
		log.Debug(ErrEmptyLinkSlice.Error())
		return ErrEmptyLinkSlice
	}

	for _, linkConf := range n.LinkDesc {

		log.Debugf("Starting link ID=[%s]\n", linkConf.linkID)

		newLink := &TCPReconn{
			lock:         sync.RWMutex{},
			status:       0,
			LinkConf:     linkConf,
			LinkActiveID: linkConf.linkID + ":" + linkConf.address,
		}

		n.LinkActives = append(n.LinkActives, &TCPReconn{})

		switch strings.ToLower(linkConf.mode) {
		case "client":
			newConn, err = n.initClientLink(newLink, linkConf)
		case "server":
			newConn, err = n.initServerLink(newLink, linkConf)
		default:
			log.Error(strings.ToLower(linkConf.mode))
			log.Warnf("%s (ID=%s)\n", ErrWrongNodeMode.Error(), linkConf.linkID)
		}

		if err != nil {
			log.Warnf("Could not connect: id=%s\n", linkConf.linkID)
			continue
		}

	}

	return nil
}

func (n *Node) initClientLink(link *TCPReconn, linkConf *LinkDesc) (*net.TCPConn, error) {

	log.Debug("Init client node")

	if n.LinkExists(link.LinkActiveID) {
		log.WithField("ID=", linkConf.linkID).Warn("Node with such name is exist yet; ignore")
	}

	raddr, err := net.ResolveTCPAddr(network, linkConf.address)
	if err != nil {
		log.WithField("ID=", linkConf.linkID).Errorf("%s error: %s", linkConf.linkID, err.Error())
		return nil, err
	}

	var conn *net.TCPConn
	for i := 0; i < 1000; i++ { // TODO: change this block when refactoring
		conn, err = net.DialTCP(network, nil, raddr)
		if err != nil {
			log.WithField("ID=", linkConf.linkID).Warnf("Connecting to %s: %s", raddr.IP, err.Error())
			//TODO: EXPONENTIAL CONNECTION
			time.Sleep(time.Second)
			continue
		} else {
			return conn, nil
		}
	}

	return nil, ErrMaxRetries
}

func (n *Node) initServerLink(link *TCPReconn, linkConf *LinkDesc) (*net.TCPConn, error) {

	log.Debug("Init server node")

	log.Debugf("ID=%s", linkConf.linkID)
	if n.LinkExists(link.LinkActiveID) {
		log.WithField("ID=", linkConf.linkID).Warnf("Node with such name is exist yet; ignore")
	}

	var conn net.Conn
	var err error
	var ln net.Listener
	for i := 0; i < 1000; i++ { // TODO: change this block when refactoring
		// listen on all interfaces
		ln, err = net.Listen(network, linkConf.address)
		if err != nil {
			log.WithField("ID=", linkConf.linkID).Errorf("Error listen: %s", err.Error())
			//TODO: EXPONENTIAL CONNECTION WITHOUT SLEEP
			time.Sleep(time.Second * 1)
			continue
		}

		// accept connection on port
		conn, err = ln.Accept()
		if err != nil {
			log.Errorf("Error accept: %s", err.Error())
			time.Sleep(time.Second * 1)
			continue
		} else {
			return conn, nil
		}
	}
	return nil, err
}

// LinkExists checks if link with such ID exists.
func (n *Node) LinkExists(linkActiveName string) bool {

	if len(n.LinkDesc) == 0 {
		log.Debugf(" %s: no connections", n.NodeID)
		return false
	}
	for _, link := range n.LinkActives {
		//log.Debug(name)
		//log.Debugf(" %v\n", key)
		if link.LinkActiveID == linkActiveName {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------------------

// // reconnect builds a new TCP connection to replace the embedded *net.TCPConn.
// //
// // TODO: keep old socket configuration (timeout, linger...).
// func (l *Link) reconnectClient() error {
// 	// set the shared status to 'reconnecting'
// 	// if it's already the case, return early, something's already trying to
// 	// reconnect
// 	if !atomic.CompareAndSwapInt32(&c.status, statusOffline, statusReconnecting) {
// 		return nil
// 	}
//
// 	raddr := c.TCPConn.RemoteAddr()
// 	conn, err := net.DialTCP(raddr.Network(), nil, raddr.(*net.TCPAddr))
// 	if err != nil {
// 		defer atomic.StoreInt32(&c.status, statusOffline)
// 		return err
// 	}
//
// 	// set new TCP socket
// 	c.TCPConn.Close()
// 	c.TCPConn = conn
//
// 	// we're back online, set shared status accordingly
// 	atomic.StoreInt32(&c.status, statusOnline)
//
// 	return nil
// }
//
// func (c *Link) reconnectServer() error {
// 	// set the shared status to 'reconnecting'
// 	// if it's already the case, return early, something's already trying to
// 	// reconnect
// 	/*	if !atomic.CompareAndSwapInt32(&c.status, statusOffline, statusReconnecting) {
// 			return nil
// 		}
//
// 		raddr := c.TCPConn.RemoteAddr()
// 		conn, err := net.DialTCP(raddr.Network(), nil, raddr.(*net.TCPAddr))
// 		if err != nil {
// 			defer atomic.StoreInt32(&c.status, statusOffline)
// 			return err
// 		}
//
// 		// set new TCP socket
// 		c.TCPConn.Close()
// 		c.TCPConn = conn
//
// 		// we're back online, set shared status accordingly
// 		atomic.StoreInt32(&c.status, statusOnline)
// 	*/
// 	return nil
// }
//
// ----------------------------------------------------------------------------

// Read wraps net.TCPConn's Read method with reconnect capabilities.
//
// It will return ErrMaxRetries if the retry limit is reached.
func (t *TCPReconn) Read(b []byte) (int, error) {
	// protect conf values (retryInterval, maxRetries...)
	c.lock.RLock()
	defer c.lock.RUnlock()

	//log.Debug("c *TCPClient Read()")

	for i := 0; i < c.maxRetries; i++ {
		if atomic.LoadInt32(&c.status) == statusOnline {
			n, err := c.TCPConn.Read(b)
			if err == nil {
				return n, err
			}

			switch e := err.(type) {
			case *net.OpError:
				if e.Err.(*os.SyscallError) == os.ErrInvalid ||
					e.Err.(syscall.Errno) == syscall.EPIPE {
					atomic.StoreInt32(&c.status, statusOffline)
				} else {
					return n, err
				}
			default:
				if err.Error() == "EOF" {
					atomic.StoreInt32(&c.status, statusOffline)
				} else {
					return n, err
				}
			}
		} else if atomic.LoadInt32(&c.status) == statusOffline {
			if err := c.reconnectClient(); err != nil {
				return -1, err
			}
		}

		// exponential backoff
		if i < (c.maxRetries - 1) {
			time.Sleep(c.retryInterval * time.Duration(math.Pow(2, float64(i))))
		}
	}

	return -1, ErrMaxRetries
}

// ReadFrom wraps net.TCPConn's ReadFrom method with reconnect capabilities.
//
// It will return ErrMaxRetries if the retry limit is reached.
func (t *TCPReconn) ReadFrom(r io.Reader) (int64, error) {
	// protect conf values (retryInterval, maxRetries...)
	c.lock.RLock()
	defer c.lock.RUnlock()

	for i := 0; i < c.maxRetries; i++ {
		if atomic.LoadInt32(&c.status) == statusOnline {
			n, err := c.TCPConn.ReadFrom(r)
			if err == nil {
				return n, err
			}
			switch e := err.(type) {
			case *net.OpError:
				if e.Err.(syscall.Errno) == syscall.ECONNRESET ||
					e.Err.(syscall.Errno) == syscall.EPIPE {
					atomic.StoreInt32(&c.status, statusOffline)
				} else {
					return n, err
				}
			default:
				if err.Error() == "EOF" {
					atomic.StoreInt32(&c.status, statusOffline)
				} else {
					return n, err
				}
			}
		} else if atomic.LoadInt32(&c.status) == statusOffline {
			switch {
			case c.Link.Mode == "client":
				if err := c.reconnectClient(); err != nil {
					return -1, err
				}

			case c.Link.Mode == "server":
				if err := c.reconnectServer(); err != nil {
					return -1, err
				}
			}

		}

		// exponential backoff
		if i < (c.maxRetries - 1) {
			time.Sleep(c.retryInterval * time.Duration(math.Pow(2, float64(i))))
		}
	}

	return -1, ErrMaxRetries
}

// Write wraps net.TCPConn's Write method with reconnect capabilities.
//
// It will return ErrMaxRetries if the retry limit is reached.
func (t *TCPReconn) Write(b []byte) (int, error) {
	// protect conf values (retryInterval, maxRetries...)
	c.lock.RLock()
	defer c.lock.RUnlock()

	log.Warn("My write() method")

	for i := 0; i < c.maxRetries; i++ {
		if atomic.LoadInt32(&c.status) == statusOnline {
			n, err := c.TCPConn.Write(b)
			if err == nil {
				return n, err
			}
			atomic.StoreInt32(&c.status, statusOffline)

		} else if atomic.LoadInt32(&c.status) == statusOffline {
			switch {
			case c.Link.Mode == "client":
				if err := c.reconnectClient(); err != nil {
					return -1, err
				}
			case c.Link.Mode == "server":
				if err := c.reconnectServer(); err != nil {
					return -1, err
				}
			}

		}

		// exponential backoff
		if i < (c.maxRetries - 1) {
			time.Sleep(c.retryInterval * time.Duration(math.Pow(2, float64(i))))
		}
	}

	return -1, ErrMaxRetries
}
