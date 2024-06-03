package commander

type Config struct {
	Prompt   string
	Commands []*Command
}

func (c *Config) ApplyDefaults() {
	if c.Prompt == "" {
		c.Prompt = "»"
	}
}
