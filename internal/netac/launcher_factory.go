package netac

type LauncherFactory interface {
	Create(config *Config) Launcher
}
