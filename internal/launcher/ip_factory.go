package launcher

import "net"

type IPFactory struct {}

func (factory *IPFactory) Create(config *Config) Launcher {
	// Determine launcher using version of config's IP.
	parsedIP := net.ParseIP(config.IP)
	if parsedIP.To4() != nil {
		return NewIPv4(config)
	}
	return NewIPv6(config)
}
