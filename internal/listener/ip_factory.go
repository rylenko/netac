package listener

import (
	"fmt"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type IPFactory struct {}

func (factory *IPFactory) Create(conn any) (listener Listener, err error) {
	switch conn.(type) {
	case *ipv4.PacketConn:
		return NewIPv4(conn.(*ipv4.PacketConn)), nil
	case *ipv6.PacketConn:
		return NewIPv6(conn.(*ipv6.PacketConn)), nil
	default:
		return nil, fmt.Errorf("unknown connection type: %T", conn)
	}
}
