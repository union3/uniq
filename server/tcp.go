package server

import (
	"fmt"
	"io"
	"net"
)

type tcpServer struct {
	ctx *context
}

func (s *tcpServer) Handle(clientConn net.Conn) {
	fmt.Println("info", "TCP: new client ", clientConn.RemoteAddr())
	buf := make([]byte, 4)
	_, err := io.ReadFull(clientConn, buf)
	if err != nil {
		fmt.Println("error", "TCP: failed to read protocol version - ", err.Error())
	}
	protocolMagic := string(buf)
	fmt.Println("info", "TCP: protocol magic ", protocolMagic, "client ", clientConn.RemoteAddr())
}
