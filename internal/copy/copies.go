package copy

import (
	"io"
	"time"
)

type Copies interface {
	DeleteExpired(ttl time.Duration)
	Print(dest io.Writer) error
	Register(newCopy Copy)
}
