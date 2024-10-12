package speaker

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/ipv6"
)

type IPv6 struct {
	idBytes []byte
	conn *ipv6.PacketConn
}

func (speaker *IPv6) SpeakForever(
		dest net.Addr, appId []byte, delay time.Duration) error {
	// Concatenate application and copy identificators to send to connection.
	buf := append(appId, speaker.idBytes...)

	for {
		if _, err := speaker.conn.WriteTo(buf, nil, dest); err != nil {
			return fmt.Errorf("failed to speak to %s: %v", dest.String(), err)
		}

		// Sleep delay before next speak.
		time.Sleep(delay)
	}
}

func NewIPv6(conn *ipv6.PacketConn) (speaker *IPv6, err error) {
	// Generate random identificator bytes.
	idBytes, err := generateRandomUUIDBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID bytes: %v", err)
	}

	speaker := &IPv6{idBytes: idBytes, conn: conn}
	return speaker, nil
}
