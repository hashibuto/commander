package commander

import (
	"fmt"

	ns "github.com/hashibuto/nilshell"
)

type Commander struct {
	Config Config

	commandMap map[string]*Command
	shell      *ns.NilShell
}

// NewCommander returns a new Commander instance
func NewCommander(config Config) (*Commander, error) {
	config.ApplyDefaults()

	commandMap := map[string]*Command{}
	for _, cmd := range config.Commands {
		if _, exists := commandMap[cmd.Name]; exists {
			return nil, fmt.Errorf("command \"%s\" is defined multiple times", cmd.Name)
		}

		err := cmd.Validate()
		if err != nil {
			return nil, err
		}

		commandMap[cmd.Name] = cmd
	}

	c := &Commander{
		Config:     config,
		commandMap: commandMap,
	}
	c.shell = ns.NewShell(config.Prompt, c.shellCompletionFunc, c.shellExecutionFunc)
	c.shell.AutoCompleteSuggestStyle = c.Config.AutoCompleteSuggestStyle

	return c, nil
}

// LocateCommand will attempt to locate a command from a series of tokens presented as arguments to the Commander.
// The method will match up to either the final subcommand, returning the remaining arguments, or to the final matching
// subcommand, returning whatever unmatched is left.
func (c *Commander) LocateCommand(tokens []string) (*Command, []string) {
	var curCommand *Command
	commandMap := c.commandMap
	lastMatchIndex := -1
	for i, token := range tokens {
		command, ok := commandMap[token]
		if !ok {
			break
		}

		lastMatchIndex = i
		curCommand = command

		commandMap = command.commandMap
		if len(commandMap) == 0 {
			break
		}
	}

	if lastMatchIndex == -1 {
		return nil, tokens
	}

	return curCommand, tokens[lastMatchIndex+1:]
}

func (c *Commander) shellCompletionFunc(beforeAndCursor string, afterCursor string, full string) []*ns.AutoComplete {
	return nil
}

func (c *Commander) shellExecutionFunc(shell *ns.NilShell, input string) {

}

func (c *Commander) Run() error {
	return c.shell.ReadUntilTerm()
}
