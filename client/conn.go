package client

import (
	"io"
	"net"
	"sync"
	"time"
)

type IdentifyResponse struct {
	AuthRequired bool `json:"auth_required"`
}

type AuthResponse struct {
	Identity string `json:"identity"`
}

type msgResponse struct {
	msg     *Message
	cmd     *Command
	success bool
	backoff bool
}

type conn struct {
	addr   net.Addr
	nc     net.Conn
	rw     io.ReadWriter
	config Config
	cli    *Client
	sync.Mutex
}

type flusher interface {
	Flush() error
}

func (c *conn) WriteCommand(cmd *Command) error {
	c.Lock()
	defer c.Unlock()
	_, err := cmd.WriteTo(c.rw)
	if err != nil {
		return err
	}
	return nil
}

func (c *conn) release() {
	c.cli.putFreeConn(c.addr, c)
}

func (c *conn) extendDeadline() {
	c.nc.SetDeadline(time.Now().Add(c.cli.netTimeout()))
}
