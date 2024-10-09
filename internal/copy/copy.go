package copy

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
)

// TODO: Add tests everywhere

type Copy struct {
	address net.Addr
	id uuid.UUID
	lastSeen time.Time
}

func (copy *Copy) Equal(other *copy) bool {
	return copy.address.String() == other.address.String() &&
		copy.id == other.id
}

func (copy *Copy) Expired(ttl time.Duration) bool {
	return time.Since(copy.lastSeen) >= ttl
}

func (copy *Copy) Extend(other *copy) {
	copy.lastSeen = other.lastSeen
}

func (copy *Copy) Print(dest io.Writer) {
	fmt.Fprintf(
		dest,
		"%s | %s | %s",
		copy.address.String(),
		copy.id.String(),
		copy.lastSeen.Format(time.TimeOnly))
}

func NewCopy(address net.Addr, id uuid.UUID, lastSeen time.Time) *Copy {
	return &Copy{
		address: address,
		id: id,
		lastSeen: lastSeen,
	}
}
