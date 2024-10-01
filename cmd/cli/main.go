package main

import (
	"log"
	"os"

	"github.com/rylenko/netac/internal/netac"
)

const help string = "$ netac <interface> <multicast-ipv4>"

func main() {
	// Validate arguments count.
	if len(os.Args) != 3 {
		log.Fatal(help)
	}

	// Extract parameters from arguments.
	ifaceName := os.Args[1]
	multicastIPv4 := os.Args[2]

	// Launch application based on the accepted parameters.
	if err := netac.Launch(ifaceName, multicastIPv4); err != nil {
		log.Fatal(err)
	}
}
