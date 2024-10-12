package printer

import (
	"fmt"
	"io"
	"time"

	"github.com/rylenko/netac/internal/copy"
)

type Delayed struct {
	delay time.Duration
}

func (printer *Delayed) PrintForever(
		copies *copy.Copies, writer io.Writer) error {
	for {
		// Print copies to the writer.
		if err := copies.Print(writer); err != nil {
			return fmt.Errorf("failed to print to the writer: %v", err)
		}

		// Sleep delay before next print.
		time.Sleep(printer.delay)
	}
}

func NewDelayed(delay time.Duration) *Delayed {
	return &Delayed{delay: delay}
}
