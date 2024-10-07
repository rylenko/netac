package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rylenko/netac/internal/netac"
)

const (
	missingRequiredParamsExitCode int = 1
)

var (
	iface *string = flag.String("iface", "eth0", "multicast interface")
	IP *string = flag.String("ip", "", "multicast IP")
)

func main() {
	flag.Parse()

	// Validate required paremeters.
	if *IP == "" {
		fmt.Fprintln(os.Stderr, "Missing required paremeters.\n")
		flag.Usage()
		os.Exit(missingRequiredParamsExitCode)
	}

	// Launch application based on the accepted parameters.
	if err := netac.Launch(context.Background(), *iface, *IP); err != nil {
		log.Fatal(err)
	}
}
