package launcher

import (
	"context"
	"fmt"
	"net"

	"github.com/rylenko/netac/internal/copy"
	"github.com/rylenko/netac/internal/listener"
	"github.com/rylenko/netac/internal/speaker"
	"github.com/rylenko/netac/internal/printer"
	"golang.org/x/net/ipv4"
)

type IPv4 struct {
	config *Config
}

func (launcher *IPv4) Launch(
		ctx context.Context,
		listenerFactory listener.Factory,
		speakerFactory speaker.Factory,
		printerImpl printer.Printer) error {
	// Try to resolve multicast address.
	addrStr := launcher.config.IP + ":" + launcher.config.Port
	multicastAddr, err := net.ResolveUDPAddr("udp4", addrStr)
	if err != nil {
		return fmt.Errorf("failed to resolve address %s: %v", addrStr, err)
	}

	// Get interface by accepted name.
	iface, err := net.InterfaceByName(launcher.config.IfaceName)
	if err != nil {
		return fmt.Errorf(
			"failed to get interface by name %s: %v", launcher.config.IfaceName, err)
	}

	// Listen packets on the accepted port.
	conn, err := getListenConfig().
		ListenPacket(ctx, "udp4", ":" + launcher.config.Port)
	if err != nil {
		return fmt.Errorf(
			"failed to listen packet on port %s: %v", launcher.config.Port, err)
	}

	// Create a new instance of low-level packet connection.
	packetConn := ipv4.NewPacketConn(conn)
	// Join multicast group using resolver address.
	if err := packetConn.JoinGroup(iface, multicastAddr); err != nil {
		return fmt.Errorf("failed to join multicast %s: %v", addrStr, err)
	}
	// Set multicast packets TTL.
	if err := packetConn.SetMulticastTTL(launcher.config.PacketTTL); err != nil {
		return fmt.Errorf(
			"failed to set packet TTL %d: %v", launcher.config.PacketTTL, err)
	}
	// Set multicast loopback.
	if err := packetConn.SetMulticastLoopback(true); err != nil {
		return fmt.Errorf("failed to set multicast loopback: %v", err)
	}

	// Create a storage of all copies.
	var copies copy.Copies

	// Get and run listener implementation.
	//
	// TODO: handle errors
	listener := listenerFactory.Create(packetConn)
	go listener.ListenForever(
		&copies, launcher.config.CopyTTL, launcher.config.AppId)

	// Get and run speaker implementation.
	//
	// TODO: handle errors
	speaker := speakerFactory.Create(packetConn)
	go speaker.SpeakForever(
		multicastAddr, launcher.config.AppId, launcher.config.SpeakDelay)

	// Print copies to standard output.
	if err := printerImpl.PrintForever(&copies, os.Stdout); err != nil {
		return fmt.Errorf("failed to print copies: %v", err)
	}
	return nil
}

func NewIPv4(config *Config) *IPv4 {
	return &IPv4{config: config}
}
