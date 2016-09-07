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
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ventu-io/slf"
)

var log = slf.WithContext("client.go")

// ----------------------------------------------------------------------------

// status enum
const (
	statusOnline       = iota
	statusOffline      = iota
	statusReconnecting = iota
)

// ----------------------------------------------------------------------------

// Link is a struct with conn options for creating a node
type LinkOpts struct {
	ID       string
	Address  string
	Mode     string
	Internal int
}

// TCPConn provides a TCP connection with auto-reconnect capabilities.
//
// It embeds a *net.TCPConn and thus implements the net.Conn interface.
//
// Use the SetMaxRetries() and SetRetryInterval() methods to configure retry
// values; otherwise they default to maxRetries=5 and retryInterval=100ms.
//
// TCPConn can be safely used from multiple goroutines.
type TCPConn struct {
	conn          *net.TCPConn
	lock          sync.RWMutex
	status        int32
	maxRetries    int
	retryInterval time.Duration
	Link          LinkOpts
}

var conns map[string]TCPConn

//Init reads links and creates nodes with their options
func Init(links []LinkOpts) {
	if links == nil && len(links) == 0 {
		log.Error("Could not init connections: empty link's slice.")
		return
	}

	conns = make(map[string]TCPConn)

	for _, link := range links {
		if link.Mode == "client" {
			initClientNode(link)
		}
		if link.Mode == "server" {
			initServerNode(link)
		}

	}
	//conns[link.ID] =
	//	}
}

func initClientNode(link LinkOpts) {
	if linkExists(link.ID) {
		log.Warn("Node with such name is exist yet; ignore")
	}
	conns[link.ID] = TCPConn{
		Link: link,
	}

	ConnectClient("tcp", link.Address, link)

}

func initServerNode(link LinkOpts) {
	if linkExists(link.ID) {
		log.Warn("Node with such name is exist yet; ignore")
	}
	conns[link.ID] = TCPConn{
		Link: link,
	}
	//conn[link.ID] = TCPConn{}
}

func linkExists(linkName string) bool {
	for name := range conns {
		if name == linkName {
			return true
		}
	}
	return false
}

// ConnectClient returns a new net.Conn.
//
// The new client connects to the remote address `raddr` on the network `network`,
// which must be "tcp", "tcp4", or "tcp6".
//
// This complements net package's Dial function.
func ConnectClient(network, addr string, links ...LinkOpts) (*TCPConn, error) {
	if links == nil{
	raddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		log.Error(err.Error())
	}
	}else{
		for _, link := range links{
			raddr, err := net.ResolveTCPAddr(network, link.Address)
		}
	}

		if  link == nil{
return DialTCP(network, nil, raddr)
		}else if len(conns) != 0{
			delete conns[link.ID]
		}	

		if 

		for _, link := range links{
			raddr, err := net.ResolveTCPAddr(network, addr)
		}
		return nil, err
	}
}

// ConnectServer creates a working server node
func ConnectServer() (net.Conn, error) {
	// listen on all interfaces
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		return nil, err
	}

	// accept connection on port
	conn, err := ln.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// DialTCP returns a new *TCPConn.
//
// The new client connects to the remote address `raddr` on the network `network`,
// which must be "tcp", "tcp4", or "tcp6".
// If `laddr` is not nil, it is used as the local address for the connection.
//
// This overrides net.TCPConn's DialTCP function.
func DialTCP(network string, laddr, raddr *net.TCPAddr) (*TCPConn, error) {

	var conn *net.TCPConn
	var err error
	for {
		conn, err = net.DialTCP(network, laddr, raddr)
		if err != nil {
			log.Warnf("Connecting to %s: %s", raddr.IP, err.Error())
			time.Sleep(time.Second)
			continue
		} else {
			break
		}
	}

	return &TCPConn{
		TCPConn: conn,

		lock:   sync.RWMutex{},
		status: 0,

		maxRetries:    10,
		retryInterval: 10 * time.Millisecond,
	}, nil
}

// ----------------------------------------------------------------------------

// SetMaxRetries sets the retry limit for the TCPClient.
//
// Assuming i is the current retry iteration, the total sleep time is
// t = retryInterval * (2^i)
//
// This function completely Lock()s the TCPClient.
func (c *TCPConn) SetMaxRetries(maxRetries int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.maxRetries = maxRetries
}

// GetMaxRetries gets the retry limit for the TCPClient.
//
// Assuming i is the current retry iteration, the total sleep time is
// t = retryInterval * (2^i)
func (c *TCPConn) GetMaxRetries() int {
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
func (c *TCPConn) SetRetryInterval(retryInterval time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.retryInterval = retryInterval
}

// GetRetryInterval gets the retry interval for the TCPClient.
//
// Assuming i is the current retry iteration, the total sleep time is
// t = retryInterval * (2^i)
func (c *TCPConn) GetRetryInterval() time.Duration {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.retryInterval
}

// ----------------------------------------------------------------------------

// reconnect builds a new TCP connection to replace the embedded *net.TCPConn.
//
// TODO: keep old socket configuration (timeout, linger...).
func (c *TCPConn) reconnect() error {
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

// ----------------------------------------------------------------------------

// Read wraps net.TCPConn's Read method with reconnect capabilities.
//
// It will return ErrMaxRetries if the retry limit is reached.
func (c *TCPConn) Read(b []byte) (int, error) {
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
			if err := c.reconnect(); err != nil {
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
func (c *TCPConn) ReadFrom(r io.Reader) (int64, error) {
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
			if err := c.reconnect(); err != nil {
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

// Write wraps net.TCPConn's Write method with reconnect capabilities.
//
// It will return ErrMaxRetries if the retry limit is reached.
func (c *TCPConn) Write(b []byte) (int, error) {
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
			if err := c.reconnect(); err != nil {
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
