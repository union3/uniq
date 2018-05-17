package server

import (
	"net"
	"sync/atomic"
	"time"
)

type protocolV1 struct {
	ctx *context
}

func (p *protocolV1) Loop(conn net.Conn) error {
	var err error
	var line []byte
	var zeroTime time.Time
	clientID := atomic.AddInt64(&p.ctx.nsqd.clientIDSequence, 1)
	client := newClientV1(clientID, conn, p.ctx)
	for {
	}
	conn.Close()
	return err
}
