package server

import (
	"net"
)

type clientV1 struct {
	ID       int64
	ctx      *context
	Hostname string
	net.Conn
}

func newClientV1(id int64, conn net.Conn, ctx *context) *clientV1 {
	var identifier string
	if conn != nil {
		identifier, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	}
	c := &clientV1{
		ID:       id,
		ctx:      ctx,
		Conn:     conn,
		Hostname: identifier,
	}
	return c
}
