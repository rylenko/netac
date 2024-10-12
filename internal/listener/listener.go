package listener

import (
	"time"

	"github.com/rylenko/netac/internal/copy"
)

type Listener interface {
	ListenForever(copies *copy.Copies, copyTTL time.Duration, appId []byte) error
}
