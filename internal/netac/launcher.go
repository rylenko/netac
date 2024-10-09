package netac

import "context"

type Launcher interface {
	Launch(ctx context.Context) error
}
