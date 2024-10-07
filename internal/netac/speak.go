package netac

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/ipv4"
)

func generateRandomUUIDBytes() (bytes []byte, err error) {
	// Generate a new copy identifactor.
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate a new identificator: %v", err)
	}

	// Marshal identificator to bytes.
	idBytes, err := id.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal identificator to bytes: %v", err)
	}
	return idBytes, nil
}

func speakForever(
		packetConn *ipv4.PacketConn,
		dest net.Addr,
		appId []byte,
		delay time.Duration) error {
	// Generate a new copy identifactor bytes.
	copyIdBytes, err := generateRandomUUIDBytes()
	if err != nil {
		return fmt.Errorf("failed to generate random UUID bytes: %v", err)
	}

	// Concatenate application and copy identificators to send to multicast group.
	buf := append(appId, copyIdBytes...)

	for {
		// Send the identity to multicast group.
		if _, err := packetConn.WriteTo(buf, nil, dest); err != nil {
			return fmt.Errorf(
				"failed to write the identity to multicast %s: %v", dest.String(), err)
		}

		// Sleep before next send.
		time.Sleep(delay)
	}
}
