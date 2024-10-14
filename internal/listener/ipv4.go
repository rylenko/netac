package listener

import (
	"bytes"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rylenko/netac/internal/copy"
	"golang.org/x/net/ipv4"
)

type IPv4 struct {
	conn *ipv4.PacketConn
}

func (listener *IPv4) ListenForever(
		copies copy.Copies, copyTTL time.Duration, appId []byte) error {
	// Buffer to read application identity and UUID bytes.
	buf := make([]byte, len(appId) + copy.ImplIdBytesLen)

	for {
		copies.DeleteExpired(copyTTL)

		// Read data to the buffer.
		_, _, src, err := listener.conn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("failed to read from connection: %v", err)
		}

		// Validate application identity.
		if !bytes.Equal(buf[:len(appId)], appId) {
			continue
		}

		// Try to parse a copy UUID.
		var copyId uuid.UUID
		copyIdBytes := buf[len(appId):len(appId) + copy.ImplIdBytesLen]
		if err := copyId.UnmarshalBinary(copyIdBytes); err != nil {
			continue
		}

		// Store a new copy in the storage.
		copy := copy.NewCopyImpl(src, copyId, time.Now())
		copies.Register(copy)
	}
}

func NewIPv4(conn *ipv4.PacketConn) *IPv4 {
	return &IPv4{conn: conn}
}
