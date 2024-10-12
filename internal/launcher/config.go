package launcher

import "time"

type Config struct {
	IfaceName string
	IP string
	Port string
	AppId string

	PacketTTL int
	CopyTTL time.Duration
	SpeakDelay time.Duration
}

func NewConfig(
		ifaceName, ip, port, appId string,
		packetTTL int,
		copyTTL, printDelay, speakDelay time.Duration) *Config {
	return &Config {
		IfaceName: ifaceName,
		IP: ip,
		Port: port,
		AppId: appId,
		PacketTTL: packetTTL,
		CopyTTL: copyTTL,
		PrintDelay: printDelay,
		SpeakDelay: speakDelay,
	}
}
