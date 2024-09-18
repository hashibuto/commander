package commander

type Config struct {
	Prompt   string
	Commands []*Command
	DumpFile string // For debugging purposes, all input will be sent to this file, if set

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
