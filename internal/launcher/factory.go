package launcher

type Factory interface {
	Create(config *Config) (launcher Launcher, err error)
}
