package turbo

import "net"

type Protocol interface {
	Loop(conn net.Conn) error
}
