package tcp

import (
	"io"
	"math"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const network = "tcp"

// ----------------------------------------------------------------------------

// status enum
const (
	statusOnline       = iota
	statusOffline      = iota
	statusReconnecting = iota
)

// ----------------------------------------------------------------------------

// LinkOpt contains options for connecting to an address.
// It is store options for creating connection for ctrtain connection.
type LinkOpt struct {
	LinkID string
	// Address is an address to use when dialing an
	// address
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
}

// Link describe one transport connection.
// It consist of Option - options,
// and Conn - an active connection with this options.
//
// Link realizes linker interface from transport library.
type Link struct {
	Option *LinkOpt
	Conn   *TCPReconn
}

// ----------------------------------------------------------------------------

// reconnect builds a new TCP connection to replace the embedded *net.TCPConn.
//
// TODO: keep old socket configuration (timeout, linger...).
func (c *Link) reconnectClient() error {
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

func (c *Link) reconnectServer() error {
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
func (c *Link) Read(b []byte) (int, error) {
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
func (c *Link) ReadFrom(r io.Reader) (int64, error) {
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
func (c *Link) Write(b []byte) (int, error) {
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
