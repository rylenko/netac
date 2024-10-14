package copy

import (
	"io"
	"net"
	"time"

	"github.com/google/uuid"
)

type Copy interface {
	Address() net.Addr
	Equal(other Copy) bool
	Expired(ttl time.Duration) bool
	Id() uuid.UUID
	LastSeen() time.Time
	Print(dest io.Writer) error
	ProlongUntil(other Copy)
}
