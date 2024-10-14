package copy

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/google/uuid"
)

// TODO: Add tests everywhere

const ImplIdBytesLen int = 16

type CopyImpl struct {
	address net.Addr
	id uuid.UUID
	lastSeen time.Time
}

func (copy *CopyImpl) Address() net.Addr {
	return copy.address
}

func (copy *CopyImpl) Equal(other Copy) bool {
	return copy.address.String() == other.Address().String() &&
		copy.id == other.Id()
}

func (copy *CopyImpl) Expired(ttl time.Duration) bool {
	return time.Since(copy.lastSeen) >= ttl
}

func (copy *CopyImpl) Id() uuid.UUID {
	return copy.id
}

func (copy *CopyImpl) LastSeen() time.Time {
	return copy.lastSeen
}

func (copy *CopyImpl) Print(dest io.Writer) error {
	_, err := fmt.Fprintf(
		dest,
		"%s | %s | %s",
		copy.address.String(),
		copy.id.String(),
		copy.lastSeen.Format(time.TimeOnly))
	return err
}

func (copy *CopyImpl) ProlongUntil(other Copy) {
	copy.lastSeen = other.LastSeen()
}

func NewCopyImpl(
		address net.Addr, id uuid.UUID, lastSeen time.Time) *CopyImpl {
	return &CopyImpl{
		address: address,
		id: id,
		lastSeen: lastSeen,
	}
}
