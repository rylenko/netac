package netac

import "net"

type IPLauncherFactory struct {}

func (factory *IPLauncherFactory) Create(config *Config) Launcher {
	// Determine launcher using version of config's IP.
	parsedIP := net.ParseIP(config.IP)
	if parsedIP.To4() != nil {
		return NewIPv4Launcher(config)
	}
	return NewIPv6Launcher(config)
}
