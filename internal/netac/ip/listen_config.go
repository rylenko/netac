package netac

import (
	"net"
	"syscall"
)

func getListenConfig() *net.ListenConfig {
	var config net.ListenConfig

	// Set controller to enable address reusing.
	config.Control = func(network, address string, conn syscall.RawConn) error {
		var err error
		err = conn.Control(func (fd uintptr) {
			err = syscall.SetsockoptInt(
				int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		})
		return err
	}

	return &config
}
