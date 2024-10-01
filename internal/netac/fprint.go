package netac

import (
	"fmt"
	"io"
	"sync"
	"time"
)

func fprint(copies sync.Map, dest io.Writer, delay time.Duration) {
	for {
		// Print each copy.
		copyNumber := 1
		copies.Range(func (addrStr, lastSeen any) bool {
			fmt.Fprintf(
				dest,
				"%d. %s [%s]\n",
				copyNumber,
				addrStr.(string),
				lastSeen.(time.Time).String())

			copyNumber++
			return true
		})

		// Sleep before next print.
		time.Sleep(delay)
	}
}
