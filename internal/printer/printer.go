package printer

import (
	"io"

	"github.com/rylenko/netac/internal/copy"
)

type Printer interface {
	PrintForever(copies copy.Copies, writer io.Writer) error
}
