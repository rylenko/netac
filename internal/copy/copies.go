package copy

import (
	"fmt"
	"io"
	"slices"
	"sync"
	"time"
)

type Copies struct {
	inner []*Copy
	mutex sync.Mutex
}

func (copies *Copies) DeleteExpired(ttl time.Duration) {
	copies.mutex.Lock()
	defer copies.mutex.Unlock()

	copies.inner = slices.DeleteFunc(copies.inner, func(copy *Copy) bool {
		return copy.Expired(ttl)
	})
}

func (copies *Copies) Print(dest io.Writer) {
	copies.mutex.Lock()
	defer copies.mutex.Unlock()

	fmt.Fprintln(dest, "\n>>> Copies:")

	for index, copy := range copies.inner {
		fmt.Fprintf(dest, "%d. ", index + 1)
		copy.Print(dest)
		fmt.Fprintln(dest, "")
	}
}

func (copies *Copies) Register(newCopy *Copy) {
	copies.mutex.Lock()
	defer copies.mutex.Unlock()

	for _, copy := range copies.inner {
		if copy.Equal(newCopy) {
			copy.Extend(newCopy)
			return
		}
	}

	copies.inner = append(copies.inner, newCopy)
}
