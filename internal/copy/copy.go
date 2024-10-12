package copy

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
)

// TODO: Add tests everywhere

const IdBytesLen int = 16

type Copy struct {
	address net.Addr
	id uuid.UUID
	lastSeen time.Time
}

func (copy *Copy) Equal(other *Copy) bool {
	return copy.address.String() == other.address.String() &&
		copy.id == other.id
}

func (copy *Copy) Expired(ttl time.Duration) bool {
	return time.Since(copy.lastSeen) >= ttl
}

func (copy *Copy) Extend(other *Copy) {
	copy.lastSeen = other.lastSeen
}

func (copy *Copy) Print(dest io.Writer) error {
	_, err := fmt.Fprintf(
		dest,
		"%s | %s | %s",
		copy.address.String(),
		copy.id.String(),
		copy.lastSeen.Format(time.TimeOnly))
	return err
}

func NewCopy(address net.Addr, id uuid.UUID, lastSeen time.Time) *Copy {
	return &Copy{
		address: address,
		id: id,
		lastSeen: lastSeen,
	}
}
