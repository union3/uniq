package server

import "net"

type UNIQ struct {
	tcpListener  net.Listener
	httpListener net.Listener
}

func New() (q *UNIQ) {
	q = &UNIQ{}
	return q
}

func (q *UNIQ) Main() {
	var err error
	ctx := &context{
		q
	}
}
