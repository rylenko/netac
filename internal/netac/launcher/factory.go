package netac

type LauncherFactory interface {
	Create(config *Config) (launcher Launcher, err error)
}
