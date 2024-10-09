package netac

import "net"

type IPLauncherFactory struct {}

func (factory *IPLauncherFactory) Create(
		config *Config) (launcher Launcher, err error) {
	// Determine launcher using version of config's IP.
	parsedIP := net.ParseIP(config.IP)
	if parsedIP.To4() {
		return NewIPv4Launcher(config), nil
	}
	if parsedIP.To6() {
		return NewIPv6Launcher(config), nil
	}
	return nil, errors.New("failed to determine IP version: %s", config.IP)
}
