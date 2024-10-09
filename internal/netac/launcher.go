package netac

type Launcher interface {
	Launch(ctx context.Context) error
}
