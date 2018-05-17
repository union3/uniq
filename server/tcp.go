package server

import (
	"fmt"
	"io"
	"net"

	"github.com/union3/uniq/internal/turbo"
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
	var proto turbo.Protocol
	switch protocolMagic {
	case "V1":
		proto = &protocolV1{
			ctx: s.ctx,
		}
	default:
		clientConn.Close()
		fmt.Println("info", "TCP: bad protocol magic ", protocolMagic, "client ", clientConn.RemoteAddr())
		return
	}
	err = proto.Loop(clientConn)
	if err != nil {
		fmt.Println("error", " client ", clientConn.RemoteAddr(), err.Error())
		return
	}
}
