package netac

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"

	"golang.org/x/net/ipv4"
)

func Launch(ctx context.Context, config *Config) error {
	// Resolve multicast address.
	multicastAddr, err := net.ResolveUDPAddr(
		"udp4", config.IP + ":" + config.Port)
	if err != nil {
		return fmt.Errorf("failed to resolve multicast %s: %v", config.IP, err)
	}

	// Get interface by accepted name.
	iface, err := net.InterfaceByName(config.IfaceName)
	if err != nil {
		return fmt.Errorf(
			"failed to get interface by name %s: %v", config.IfaceName, err)
	}

	// Listen packets on the constant port.
	conn, err := getListenConfig().ListenPacket(ctx, "udp4", ":" + config.Port)
	if err != nil {
		return fmt.Errorf("failed to listen packet on port %s: %v", config.Port, err)
	}
	defer conn.Close()

	// Create a new instance of packet connection.
	packetConn := ipv4.NewPacketConn(conn)
	// Join multicast group using resolver address.
	if err := packetConn.JoinGroup(iface, multicastAddr); err != nil {
		return fmt.Errorf(
			"failed to join multicast %s: %v", multicastAddr.String(), err)
	}
	// Set multicast TTL for outcoming packets.
	if err := packetConn.SetMulticastTTL(config.PacketTTL); err != nil {
		return fmt.Errorf("failed to set packet TTL %d: %v", config.PacketTTL, err)
	}
	// Set multicast loopback.
	if err := packetConn.SetMulticastLoopback(true); err != nil {
		return fmt.Errorf("failed to set multicast loopback: %v", err)
	}

	// Storage of all copies. Keys are address strings, values are copy slices.
	var copies Copies
	// Listen incoming packets.
	//
	// TODO: handle error, use config struct.
	go listenForever(&copies, packetConn, config.CopyTTL, config.AppId)
	// Print copies to writer.
	//
	// TODO: use config struct.
	go printForever(&copies, os.Stdout, config.PrintDelay)
	// Speak to multicast address.
	err = speakForever(packetConn, multicastAddr, config.AppId, config.SpeakDelay)
	if err != nil {
		return fmt.Errorf("failed to speak to %s: %v", multicastAddr.String(), err)
	}
	return nil
}

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
