package netac

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
)

func listen(
		copies sync.Map,
		packetConn *ipv4.PacketConn,
		copyTTL time.Duration,
		identity []byte) error {
	// Buffer to read identity bytes.
	buf := make([]byte, len(identity))

	for {
		// Delete expired copies.
		copies.Range(func (addrStr, lastSeen any) bool {
			if time.Since(lastSeen.(time.Time)) >= copyTTL {
				copies.Delete(addrStr.(string))
			}
			return true
		})

		// Read data to the buffer.
		n, _, src, err := packetConn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("failed to read from connection: %v", err)
		}

		// TODO: logger.
		fmt.Printf("[Listener] Readed %d bytes.\n", n)

		// Register address if  readed data is equal to identity bytes.
		if bytes.Equal(buf, identity) {
			copies.Store(src.String(), time.Now())
		}
	}
}
