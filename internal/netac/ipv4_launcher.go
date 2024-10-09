package netac

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rylenko/netac/internal/copy"
	"golang.org/x/net/ipv4"
)

type IPv4Launcher struct {
	config *Config
}

func (launcher *IPv4Launcher) Launch(ctx context.Context) error {
	// Resolve multicast address.
	multicastAddr, err := net.ResolveUDPAddr(
		"udp4", launcher.config.IP + ":" + launcher.config.Port)
	if err != nil {
		return fmt.Errorf(
			"failed to resolve multicast %s: %v", launcher.config.IP, err)
	}
	// Get interface by accepted name.
	iface, err := net.InterfaceByName(launcher.config.IfaceName)
	if err != nil {
		return fmt.Errorf(
			"failed to get interface by name %s: %v", launcher.config.IfaceName, err)
	}
	// Listen packets on the constant port.
	conn, err := getListenConfig().
		ListenPacket(ctx, "udp4", ":" + launcher.config.Port)
	if err != nil {
		return fmt.Errorf(
			"failed to listen packet on port %s: %v", launcher.config.Port, err)
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
	if err := packetConn.SetMulticastTTL(launcher.config.PacketTTL); err != nil {
		return fmt.Errorf(
			"failed to set packet TTL %d: %v", launcher.config.PacketTTL, err)
	}
	// Set multicast loopback.
	if err := packetConn.SetMulticastLoopback(true); err != nil {
		return fmt.Errorf("failed to set multicast loopback: %v", err)
	}

	// Storage of all copies. Keys are address strings, values are copy slices.
	var copies copy.Copies
	// Listen incoming packets.
	//
	// TODO: handle error, use config struct.
	go launcher.listenForever(&copies, packetConn)
	// Print copies to writer.
	//
	// TODO: use config struct.
	go launcher.printForever(&copies, os.Stdout)
	// Speak to multicast address.
	if err := launcher.speakForever(packetConn, multicastAddr); err != nil {
		return fmt.Errorf("failed to speak to %s: %v", multicastAddr.String(), err)
	}
	return nil
}

func (launcher *IPv4Launcher) listenForever(
		copies *copy.Copies, packetConn *ipv4.PacketConn) error {
	appIdBytes := []byte(launcher.config.AppId)

	// Buffer to read identity and UUID bytes.
	buf := make([]byte, len(appIdBytes) + uuidBytesLen)

	for {
		copies.DeleteExpired(launcher.config.CopyTTL)

		// Read data to the buffer.
		_, _, src, err := packetConn.ReadFrom(buf)
		if err != nil {
			return fmt.Errorf("failed to read from connection: %v", err)
		}

		// Validate application identificator.
		if !bytes.Equal(buf[:len(appIdBytes)], appIdBytes) {
			continue
		}

		// Try to parse copy identificator.
		var copyId uuid.UUID
		copyIdBytes := buf[len(appIdBytes):len(appIdBytes) + uuidBytesLen]
		err = copyId.UnmarshalBinary(copyIdBytes)
		if err != nil {
			continue
		}

		// Store a new copy in the storage.
		copy := copy.NewCopy(src, copyId, time.Now())
		copies.Register(copy)
	}
}

func (launcher *IPv4Launcher) printForever(
		copies *copy.Copies, writer io.Writer) {
	for {
		copies.Print(writer)
		time.Sleep(launcher.config.PrintDelay)
	}
}

func (launcher *IPv4Launcher) speakForever(
		packetConn *ipv4.PacketConn, dest net.Addr) error {
	// Generate a new copy identifactor bytes.
	copyIdBytes, err := generateRandomUUIDBytes()
	if err != nil {
		return fmt.Errorf("failed to generate random UUID bytes: %v", err)
	}

	// Concatenate application and copy identificators to send to multicast group.
	buf := append([]byte(launcher.config.AppId), copyIdBytes...)

	for {
		// Send the identity to multicast group.
		if _, err := packetConn.WriteTo(buf, nil, dest); err != nil {
			return fmt.Errorf(
				"failed to write the identity to multicast %s: %v", dest.String(), err)
		}

		// Sleep before next send.
		time.Sleep(launcher.config.SpeakDelay)
	}
}

func NewIPv4Launcher(config *Config) *IPv4Launcher {
	return &IPv4Launcher{config: config}
}
