package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/union3/uniq/internal/turbo"
)

type UNIQ struct {
	tcpListener      net.Listener
	httpListener     net.Listener
	wg               sync.WaitGroup
	clientIDSequence int64
}

func New() (q *UNIQ) {
	q = &UNIQ{
		clientIDSequence: 0,
	}
	return q
}

func (q *UNIQ) Main() {
	var err error
	ctx := &context{
		q,
	}
	q.tcpListener, err = net.Listen("tcp", ":1080")
	if err != nil {
		fmt.Println("error", "listen :1080 failed", err.Error())
	}
	tcpServer := &tcpServer{
		ctx: ctx,
	}
	q.wg.Add(1)
	go func() {
		turbo.TCPServer(q.tcpListener, tcpServer)
		q.wg.Done()
	}()
}

func (q *UNIQ) Exit() {
	fmt.Println("info", "UNIQ exit")
	if q.tcpListener != nil {
		q.tcpListener.Close()
	}
}
