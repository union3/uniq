package client

import (
	"bufio"
	"net"
	"sync"
	"time"
)

const (
	DefaultTimeout      = 3 * time.Second
	DefaultMaxIdleConns = 3
)

type Client struct {
	id           int64
	addr         net.Addr
	wg           sync.WaitGroup
	timeout      time.Duration
	freeConn     map[string][]*conn
	maxIdleConns int
	sync.Mutex
}

func New() (cli *Client) {

}

func (c *Client) dial(addr net.Addr) (net.Conn, error) {
	nc, err := net.DialTimeout(addr.Network(), addr.String(), c.netTimeout())
	if err != nil {
		return nc, nil
	}
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		return nil, &ConnectTimeoutError{addr}
	}
	return nil, err
}

func (c *Client) netTimeout() time.Duration {
	if c.timeout != 0 {
		return c.timeout
	}
	return DefaultTimeout
}

func (c *Client) getConn(addr net.Addr) (*conn, error) {
	cn, ok := c.getFreeConn(addr)
	if ok {
		cn.extendDeadline()
		return cn, nil
	}
	nc, err := c.dial(addr)
	if err != nil {
		return nil, err
	}
	cn = &conn{
		nc:   nc,
		addr: addr,
		rw:   bufio.NewReadWriter(bufio.NewReader(nc), bufio.NewWriter(nc)),
	}
	cn.extendDeadline()
	return cn, nil
}

func (c *Client) getFreeConn(addr net.Addr) (cn *conn, ok bool) {
	c.Lock()
	defer c.Unlock()
	if c.freeConn == nil {
		return nil, false
	}
	freeList, ok := c.freeConn[addr.String()]
	if !ok || len(freeList) == 0 {
		return nil, false
	}
	cn = freeList[len(freeList)-1]
	c.freeConn[addr.String()] = freeList[:len(freeList)-1]
	return cn, true
}

func (c *Client) putFreeConn(addr net.Addr, cn *conn) {
	c.Lock()
	defer c.Unlock()
	if c.freeConn == nil {
		c.freeConn = make(map[string][]*conn)
	}
	freeList := c.freeConn[addr.String()]
	if len(freeList) >= c.maxIdleConns {
		cn.nc.Close()
	} else {
		c.freeConn[addr.String()] = append(freeList, cn)
	}
}
