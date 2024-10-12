package launcher

import (
	"context"

	"github.com/rylenko/netac/internal/listener"
	"github.com/rylenko/netac/internal/printer"
	"github.com/rylenko/netac/internal/speaker"
)

type Launcher interface {
	Launch(
		ctx context.Context,
		listenerFactory listener.Factory,
		speakerFactory speaker.Factory,
		printerImpl printer.Printer) error
}
