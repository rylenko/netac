package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rylenko/netac/internal/netac"
)

const (
	missingRequiredParamsExitCode int = 1
)

var (
	iface *string = flag.String("iface", "eth0", "multicast interface")
	ip *string = flag.String("ip", "", "multicast IP")
	port *string = flag.String("port", "9999", "multicast port")
	appId *string = flag.String(
		"appid", "12345", "application id for identifying copies on the network")

	packetTTL = flag.Int(
		"packetTTL", 2, "TTL of multicast packets")
	copyTTL = flag.Duration(
		"copyTTL", 10 * time.Second, "TTL of registered copies")
	printDelay = flag.Duration(
		"printDelay", 4 * time.Second, "delay of copies printing")
	speakDelay = flag.Duration(
		"speakDelay", 1 * time.Second, "delay of current copy speaking")
)

func main() {
	flag.Parse()

	// Validate required paremeters.
	if *ip == "" {
		fmt.Fprintln(os.Stderr, "Missing required paremeters.\n")
		flag.Usage()
		os.Exit(missingRequiredParamsExitCode)
	}

	// Build config using parsed arguments.
	config := netac.NewConfig(
		*iface, *ip, *port, *appId, *packetTTL, *copyTTL, *printDelay, *speakDelay)

	// Launch application based on the built config.
	var factory netac.IPLauncherFactory
	launcher := factory.Create(config)
	if err := launcher.Launch(context.Background()); err != nil {
		log.Fatal(err)
	}
}
