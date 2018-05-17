package server

import (
	"bufio"
	"net"
	"time"
)

const defaultBufferSize = 16 * 1024

type clientV1 struct {
	ID          int64
	ctx         *context
	Hostname    string
	R           *bufio.Reader
	W           *bufio.Writer
	ConnectTime int64
	net.Conn
}

func newClientV1(id int64, conn net.Conn, ctx *context) *clientV1 {
	var identifier string
	if conn != nil {
		identifier, _, _ = net.SplitHostPort(conn.RemoteAddr().String())
	}
	c := &clientV1{
		ID:          id,
		ctx:         ctx,
		Conn:        conn,
		R:           bufio.NewReaderSize(conn, defaultBufferSize),
		W:           bufio.NewWriterSize(conn, defaultBufferSize),
		ConnectTime: time.Now().Unix(),
		Hostname:    identifier,
	}
	return c
}
