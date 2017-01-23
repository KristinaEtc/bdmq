package tcp

import (
	"net"
	"time"
)

// Node is a struct which conbine connections with common options.
//
// TODO:Noder realizes linker interface from transport library. (?)
// (implements method GetStatus() for group of links).
type Node struct {
	NodeID        string
	MaxRetryCount int
	Links         []*Link
}

//RunConn reads links and creates nodes with their options
func (n *Node) RunConn() {
	//func (g *Global) RunConn() (map[string]net.Conn, error) {

	/*
		log.Debug("func Init")
		//return nil, fmt.Errorf("test")

		var err error
		var newConn net.Conn

		//log.Debugf("Init %v config\n", links)
		if links == nil && len(links) == 0 {
			log.Debug(ErrEmptyNodeSlice.Error())
			return nil, ErrEmptyNodeSlice
		}

		for num, link := range links {

			log.Debugf("%d\n", num)

			newNode := &TCPReconn{
				lock:          sync.RWMutex{},
				status:        0,
				Link:          link,
				maxRetries:    link.MaxRetries,
				retryInterval: time.Duration(link.Internal) * time.Millisecond,
			}

			switch strings.ToLower(link.Mode) {
			case "client":
				newConn, err = t.initClientNode(newNode)
			case "server":
				newConn, err = t.initServerNode(newNode)
			default:
				log.Error(strings.ToLower(link.Mode))
				log.Warnf("%s (ID=%s)\n", ErrWrongNodeMode.Error(), link.ID)
			}

			if err != nil {
				log.Warnf("Could not connect: id=%s\n", link.ID)
				continue
			}
			log.Warnf("ID=%s", newNode.Link.ID)
			t.Conns[newNode.Link.ID] = newConn
		}
		return t.Conns, nil
	*/

}

func (t *Node) initClientNode(c *Node) (*net.TCPConn, error) {

	log.Debug("Init client node")

	if t.LinkExists(c.Link.ID) {
		log.Warn("Node with such name is exist yet; ignore")
	}

	raddr, err := net.ResolveTCPAddr(network, c.Link.Address)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var conn *net.TCPConn
	for i := 0; i < c.maxRetries; i++ {
		conn, err = net.DialTCP(network, nil, raddr)
		if err != nil {
			log.Warnf("Connecting to %s: %s", raddr.IP, err.Error())
			time.Sleep(time.Second)
			continue
		} else {
			return conn, nil
		}
	}

	return nil, ErrMaxRetries
}

func (t *Node) initServerNode(c *TCPReconn) (net.Conn, error) {

	log.Debug("Init server node")

	log.Debugf("ID=%s", c.Link.ID)
	if t.LinkExists(c.Link.ID) {
		log.Warn("Node with such name is exist yet; ignore")
	}

	var conn net.Conn
	var err error
	var ln net.Listener
	for i := 0; i < c.maxRetries; i++ {
		// listen on all interfaces
		ln, err = net.Listen(network, c.Link.Address)
		if err != nil {
			log.Errorf("error listen: %s", err.Error())
			log.Error(err.Error())
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
func (t *Node) LinkExists(linkName string) bool {

	if len(t.Conns) == 0 {
		log.Debug("no connections")
		return false
	}
	for name := range t.Conns {
		//log.Debug(name)
		//log.Debugf(" %v\n", key)
		if name == linkName {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------------------

// SetMaxRetries sets the retry limit for the TCPClient.
//
// Assuming i is the current retry iteration, the total sleep time is
// t = retryInterval * (2^i)
//
// This function completely Lock()s the TCPClient.
func (c *Node) SetMaxRetries(maxRetries int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.maxRetries = maxRetries
}

// GetMaxRetries gets the retry limit for the TCPClient.
//
// Assuming i is the current retry iteration, the total sleep time is
// t = retryInterval * (2^i)
func (c *Node) GetMaxRetries() int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.maxRetries
}

// SetRetryInterval sets the retry interval for the TCPClient.
//
// Assuming i is the current retry iteration, the total sleep time is
// t = retryInterval * (2^i)
//
// This function completely Lock()s the TCPClient.
func (c *Node) SetRetryInterval(retryInterval time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.retryInterval = retryInterval
}

// GetRetryInterval gets the retry interval for the TCPClient.
//
// Assuming i is the current retry iteration, the total sleep time is
// t = retryInterval * (2^i)
func (c *Node) GetRetryInterval() time.Duration {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.retryInterval
}
