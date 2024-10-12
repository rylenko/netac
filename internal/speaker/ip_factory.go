package speaker

import (
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type IPFactory struct {}

func (factory *IPFactory) Create(conn any) (speaker Speaker, err error) {
	switch conn.(type) {
	case *ipv4.PacketConn:
		return NewIPv4(conn)
	case *ipv6.PacketConn:
		return NewIPv6(conn)
	default:
		return fmt.Errorf("unknown connection type: %T", conn)
	}
}
