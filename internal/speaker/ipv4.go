package speaker

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

type IPv4 struct {
	idBytes []byte
	conn *ipv4.PacketConn
}

func (speaker *IPv4) SpeakForever(
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

func NewIPv4(conn *ipv4.PacketConn) (speaker *IPv4, err error) {
	// Generate random identificator bytes.
	idBytes, err := generateRandomUUIDBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID bytes: %v", err)
	}

	speaker = &IPv4{idBytes: idBytes, conn: conn}
	return speaker, nil
}
