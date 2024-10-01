package netac

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

func speak(
		packetConn *ipv4.PacketConn,
		dest net.Addr,
		identity []byte,
		delay time.Duration) error {
	for {
		// Send the identity to multicast group.
		if _, err := packetConn.WriteTo(identity, nil, dest); err != nil {
			return fmt.Errorf(
				"failed to write the identity to multicast %s: %v", dest.String(), err)
		}

		// TODO: logger.
		fmt.Println("[Speaker] Identity sent.")

		// Sleep before next send.
		time.Sleep(delay)
	}
}
