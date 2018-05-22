package client

import (
	"fmt"
	"net"
)

type ErrIdentify struct {
	Reason string
}

func (e ErrIdentify) Error() string {
	return fmt.Sprintf("failed to IDENTIFY - %s", e.Reason)
}

type ConnectTimeoutError struct {
	Addr net.Addr
}

func (cte *ConnectTimeoutError) Error() string {
	return "memcache: connect timeout to " + cte.Addr.String()
}
