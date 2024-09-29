package commander

type Config struct {
	PromptFunc func() string
	Commands   []*Command
	DumpFile   string // For debugging purposes, all input will be sent to this file, if set
}
