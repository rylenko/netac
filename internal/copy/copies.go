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

func (copies *Copies) Print(dest io.Writer) error {
	copies.mutex.Lock()
	defer copies.mutex.Unlock()

	if _, err := fmt.Fprintln(dest, "\n>>> Copies:"); err != nil {
		return fmt.Errorf("failed to print copies header: %v", err)
	}

	for index, copy := range copies.inner {
		if _, err := fmt.Fprintf(dest, "%d. ", index + 1); err != nil {
			return fmt.Errorf("failed to print copy number: %v", err)
		}
		if err := copy.Print(dest); err != nil {
			return fmt.Errorf("failed to print copy: %v", err)
		}
		if _, err := fmt.Fprintln(dest, ""); err != nil {
			return fmt.Errorf("failed to print newline: %v", err)
		}
	}
	return nil
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
