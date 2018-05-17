package server

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

type protocolV1 struct {
	ctx *context
}

func (p *protocolV1) Loop(conn net.Conn) error {
	var err error
	var line []byte
	clientID := atomic.AddInt64(&p.ctx.uniq.clientIDSequence, 1)
	client := newClientV1(clientID, conn, p.ctx)
	for {
		line, err = client.R.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = fmt.Errorf("failed to read command - %s", err)
			}
			break
		}
		fmt.Println("line", line)
	}
	conn.Close()
	return err
}
