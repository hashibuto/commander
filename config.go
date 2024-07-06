package commander

type Config struct {
	Prompt   string
	Commands []*Command

	AutoCompleteSuggestStyle string
}

func (c *Config) ApplyDefaults() {
	if c.Prompt == "" {
		c.Prompt = "» "
	}

	if c.AutoCompleteSuggestStyle == "" {
		c.AutoCompleteSuggestStyle = "\033[32m"
	}
}
