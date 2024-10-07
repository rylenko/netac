package netac

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"golang.org/x/net/ipv4"
)

const (
	// TODO: create config structure to change this parameters
	copyTTL time.Duration = 10 * time.Second
	multicastTTL int = 2
	portStr string = "9999"
	printDelay time.Duration = 4 * time.Second
	speakDelay time.Duration = 1 * time.Second
)

// Bytes to send over UDP to identify application copies.
var appId []byte = []byte{1, 2, 3, 4, 5, 6, 7, 8}

func Launch(ctx context.Context, ifaceName, multicastIPv4 string) error {
	// Resolve multicast address.
	multicastAddr, err := net.ResolveUDPAddr(
		"udp4", multicastIPv4 + ":" + portStr)
	if err != nil {
		return fmt.Errorf("failed to resolve multicast %s: %v", multicastIPv4, err)
	}

	// Get interface by accepted name.
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return fmt.Errorf("failed to get interface by name %s: %v", ifaceName, err)
	}

	// Listen packets on the constant port.
	conn, err := getListenConfig().ListenPacket(ctx, "udp4", ":" + portStr)
	if err != nil {
		return fmt.Errorf("failed to listen packet on port %s: %v", portStr, err)
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
	if err := packetConn.SetMulticastTTL(multicastTTL); err != nil {
		return fmt.Errorf("failed to set multicast TTL %d: %v", multicastTTL, err)
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
	go listenForever(&copies, packetConn, copyTTL, appId)
	// Print copies to writer.
	//
	// TODO: use config struct.
	go printForever(&copies, os.Stdout, printDelay)
	// Speak to multicast address.
	err = speakForever(packetConn, multicastAddr, appId, speakDelay)
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
