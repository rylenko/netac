package speaker

import (
	"net"
	"time"
)

type Speaker interface {
	SpeakForever(dest net.Addr, appId []byte, delay time.Duration) error
}
