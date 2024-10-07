package netac

import (
	"io"
	"time"
)

func printForever(copies *Copies, writer io.Writer, delay time.Duration) {
	for {
		copies.Print(writer)
		time.Sleep(delay)
	}
}
