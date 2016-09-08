// Copyright © 2015 Clement 'cmc' Rey <cr.rey.clement@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tcprec

import (
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

var log = slf.WithContext("client.go")

const network = "tcp"

// ----------------------------------------------------------------------------

// status enum
const (
	statusOnline       = iota
	statusOffline      = iota
	statusReconnecting = iota
)

// ----------------------------------------------------------------------------

// LinkOpts is a struct with conn options for creating a node
type LinkOpts struct {
	ID         string
	Address    string
	Mode       string
	Internal   int
	MaxRetries int
}

// TCPReconn provides a TCP connection with auto-reconnect capabilities.
type TCPReconn struct {
	TCPConn       *net.TCPConn
	lock          sync.RWMutex
	status        int32
	maxRetries    int
	retryInterval time.Duration
	Link          LinkOpts
}

var conns map[string]net.Conn

//Init reads links and creates nodes with their options
func Init(links []LinkOpts) (map[string]net.Conn, error) {

	var err error
	var newConn net.Conn

	conns = make(map[string]net.Conn, 0)

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
			newConn, err = initClientNode(newNode)
		case "server":
			newConn, err = initServerNode(newNode)
		default:
			log.Error(strings.ToLower(link.Mode))
			log.Warnf("%s (ID=%s)\n", ErrWrongNodeMode.Error(), link.ID)
		}

		if err != nil {
			log.Warnf("Could not connect: id=%s\n", link.ID)
			continue
		}
		conns[newNode.Link.ID] = newConn
	}
	return conns, nil
}

func initClientNode(c *TCPReconn) (*net.TCPConn, error) {

	log.Debug("Init client node")

	if LinkExists(c.Link.ID) {
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

func initServerNode(c *TCPReconn) (net.Conn, error) {

	log.Debug("Init server node")

	if LinkExists(c.Link.ID) {
		log.Warn("Node with such name is exist yet; ignore")
	}

	var conn net.Conn
	var err error
	var ln net.Listener
	for i := 0; i < c.maxRetries; i++ {
		// listen on all interfaces
		ln, err = net.Listen(network, c.Link.Address)
		if err != nil {
			log.Error(err.Error())
			time.Sleep(time.Second * 1)
			continue
		}

		// accept connection on port
		conn, err = ln.Accept()
		if err != nil {
			log.Error(err.Error())
			time.Sleep(time.Second * 1)
			continue
		} else {
			return conn, nil
		}
	}
	return nil, err
}

func LinkExists(linkName string) bool {

	if len(conns) == 0 {
		log.Debug("no connections")
		return false
	}
	for name := range conns {
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
func (c *TCPReconn) SetMaxRetries(maxRetries int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.maxRetries = maxRetries
}

// GetMaxRetries gets the retry limit for the TCPClient.
//
// Assuming i is the current retry iteration, the total sleep time is
// t = retryInterval * (2^i)
func (c *TCPReconn) GetMaxRetries() int {
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
func (c *TCPReconn) SetRetryInterval(retryInterval time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.retryInterval = retryInterval
}

// GetRetryInterval gets the retry interval for the TCPClient.
//
// Assuming i is the current retry iteration, the total sleep time is
// t = retryInterval * (2^i)
func (c *TCPReconn) GetRetryInterval() time.Duration {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.retryInterval
}

// ----------------------------------------------------------------------------

// reconnect builds a new TCP connection to replace the embedded *net.TCPConn.
//
// TODO: keep old socket configuration (timeout, linger...).
func (c *TCPReconn) reconnectClient() error {
	// set the shared status to 'reconnecting'
	// if it's already the case, return early, something's already trying to
	// reconnect
	if !atomic.CompareAndSwapInt32(&c.status, statusOffline, statusReconnecting) {
		return nil
	}

	raddr := c.TCPConn.RemoteAddr()
	conn, err := net.DialTCP(raddr.Network(), nil, raddr.(*net.TCPAddr))
	if err != nil {
		defer atomic.StoreInt32(&c.status, statusOffline)
		return err
	}

	// set new TCP socket
	c.TCPConn.Close()
	c.TCPConn = conn

	// we're back online, set shared status accordingly
	atomic.StoreInt32(&c.status, statusOnline)

	return nil
}

func (c *TCPReconn) reconnectServer() error {
	// set the shared status to 'reconnecting'
	// if it's already the case, return early, something's already trying to
	// reconnect
	/*	if !atomic.CompareAndSwapInt32(&c.status, statusOffline, statusReconnecting) {
			return nil
		}

		raddr := c.TCPConn.RemoteAddr()
		conn, err := net.DialTCP(raddr.Network(), nil, raddr.(*net.TCPAddr))
		if err != nil {
			defer atomic.StoreInt32(&c.status, statusOffline)
			return err
		}

		// set new TCP socket
		c.TCPConn.Close()
		c.TCPConn = conn

		// we're back online, set shared status accordingly
		atomic.StoreInt32(&c.status, statusOnline)
	*/
	return nil
}

// ----------------------------------------------------------------------------

// Read wraps net.TCPConn's Read method with reconnect capabilities.
//
// It will return ErrMaxRetries if the retry limit is reached.
func (c *TCPReconn) Read(b []byte) (int, error) {
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
func (c *TCPReconn) ReadFrom(r io.Reader) (int64, error) {
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
func (c *TCPReconn) Write(b []byte) (int, error) {
	// protect conf values (retryInterval, maxRetries...)
	c.lock.RLock()
	defer c.lock.RUnlock()

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
