package turbo

import (
	"fmt"
	"net"
	"runtime"
	"strings"
)

type TCPHandler interface {
	Handle(net.Conn)
}

func TCPServer(listener net.Listener, handler TCPHandler) {
	fmt.Printf("info", "TCP: listening on %s", listener.Addr())
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				fmt.Printf("error", "temporary Accept() failure - %s", err)
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				fmt.Printf("error", "listener.Accept() - %s", err)
			}
			break
		}
		go handler.Handle(clientConn)
	}
	fmt.Printf("info", "TCP: closing %s", listener.Addr())
}
