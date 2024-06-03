package commander

type Commander struct {
	Config Config
}

// NewCommander returns a new Commander instance
func NewCommander(config Config) *Commander {
	config.ApplyDefaults()

	return &Commander{}
}
