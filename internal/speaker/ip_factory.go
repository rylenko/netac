package speaker

import (
	"fmt"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type IPFactory struct {}

func (factory *IPFactory) Create(conn any) (speaker Speaker, err error) {
	switch conn.(type) {
	case *ipv4.PacketConn:
		speaker, err := NewIPv4(conn.(*ipv4.PacketConn))
		if err != nil {
			return nil, fmt.Errorf("failed to instance IPv4 speaker: %v", err)
		}
		return speaker, nil
	case *ipv6.PacketConn:
		speaker, err := NewIPv6(conn.(*ipv6.PacketConn))
		if err != nil {
			return nil, fmt.Errorf("failed to instance IPv4 speaker: %v", err)
		}
		return speaker, nil
	default:
		return nil, fmt.Errorf("unknown connection type: %T", conn)
	}
}
