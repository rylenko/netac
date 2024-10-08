package netac

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/ipv4"
)

func listenForever(
		copies *Copies,
		packetConn *ipv4.PacketConn,
		copyTTL time.Duration,
		appId string) error {
	// Buffer to read identity and UUID bytes.
	buf := make([]byte, len(appId) + CopyIdBytesLen)
	appIdBytes := []byte(appId)

	for {
		copies.DeleteExpired(copyTTL)

		// Read data to the buffer.
		_, _, src, err := packetConn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("failed to read from connection: %v", err)
		}

		// Validate application identificator.
		if !bytes.Equal(buf[:len(appIdBytes)], appIdBytes) {
			continue
		}

		// Try to parse copy identificator.
		var copyId uuid.UUID
		copyIdBytes := buf[len(appIdBytes):len(appIdBytes) + CopyIdBytesLen]
		err = copyId.UnmarshalBinary(copyIdBytes)
		if err != nil {
			continue
		}

		// Store a new copy in the storage.
		newCopy := NewCopy(src, copyId, time.Now())
		copies.Register(newCopy)
	}
}
