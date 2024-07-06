package commander

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	ns "github.com/hashibuto/nilshell"
	"github.com/hashibuto/nilshell/pkg/term"
)

const (
	COMMAND_PADDING = 30
)

type Commander struct {
	Config Config

	commandMap map[string]*Command
	shell      *ns.NilShell
}

type BoundExec struct {
	IsCapturingOutput bool
	Command           *Command
	ArgMap            map[string]any
}

type Capture struct {
	Buffer bytes.Buffer
	Error  error
}

// NewCommander returns a new Commander instance
func NewCommander(config Config) (*Commander, error) {
	config.ApplyDefaults()
	config.Commands = append(config.Commands, HelpCommand, GrepCommand, ClearCommand)

	c := &Commander{
		Config: config,
	}

	commandMap := map[string]*Command{}
	for _, cmd := range config.Commands {
		if _, exists := commandMap[cmd.Name]; exists {
			return nil, fmt.Errorf("command \"%s\" is defined multiple times", cmd.Name)
		}

		err := cmd.Validate(nil)
		if err != nil {
			return nil, err
		}

		commandMap[cmd.Name] = cmd

		// Give everyone a reference to the commander, in order to carry out top level operations if necessary
		cmd.Commander = c
		for _, sub := range cmd.SubCommands {
			sub.Commander = c
		}
	}

	c.commandMap = commandMap
	c.shell = ns.NewShell(config.Prompt, c.shellCompletionFunc, c.shellExecutionFunc)
	c.shell.AutoCompleteSuggestStyle = c.Config.AutoCompleteSuggestStyle

	return c, nil
}

// LocateCommand will attempt to locate a command from a series of tokens presented as arguments to the Commander.
// The method will match up to either the final subcommand, returning the remaining arguments, or to the final matching
// subcommand, returning whatever unmatched is left.
func (c *Commander) LocateCommand(tokens []string) (*Command, []*Flag, []string) {
	flags := []*Flag{}

	var curCommand *Command
	commandMap := c.commandMap
	lastMatchIndex := -1
	for i, token := range tokens {
		// Grab parent command flags
		if curCommand != nil {
			flags = append(flags, curCommand.Flags...)
		}

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
		return nil, nil, tokens
	}

	return curCommand, flags, tokens[lastMatchIndex+1:]
}

// shellCompletionFunc is invoked when the user engages the tab completion feature of the shell.  this attempts to return
// suggestions for completion.
func (c *Commander) shellCompletionFunc(beforeAndCursor string, afterCursor string, full string) []*ns.AutoComplete {
	autoComplete := []*ns.AutoComplete{}
	tokenGroups := Tokenize(beforeAndCursor)
	if len(tokenGroups) == 0 {
		return nil
	}

	// we are only concerned with the last token group
	tokens := tokenGroups[len(tokenGroups)-1].Tokens

	if strings.HasSuffix(beforeAndCursor, " ") {
		// this assists in finding the "next" thing, when there's no non-whitespace input
		tokens = append(tokens, "")
	}
	command, parentFlags, remaining := c.LocateCommand(tokens)
	if len(remaining) == 0 {
		return nil
	}

	if command == nil {
		// attempt to lookup the command by partial match
		for _, lookupCmd := range c.Config.Commands {
			if strings.HasPrefix(lookupCmd.Name, remaining[0]) {
				autoComplete = append(autoComplete, &ns.AutoComplete{
					Value:   lookupCmd.Name,
					Display: lookupCmd.Name,
				})
			}
		}

		return autoComplete
	}

	suggestions := command.Suggest(remaining, parentFlags)
	autoComplete = append(autoComplete, suggestions...)

	return autoComplete
}

func (c *Commander) shellExecutionFunc(shell *ns.NilShell, input string) {
	tokenGroups := Tokenize(input)
	if len(tokenGroups) == 0 {
		return
	}

	// make sure the groups make sense first
	execSequence := []*BoundExec{}
	for i, tokenGroup := range tokenGroups {

		if tokenGroup.FlowControl == FLOW_CONTROL_REDIRECT {
			if i == 0 {
				Errorln("nothing to redirect")
				return
			}

			if i != len(tokenGroups)-1 {
				Errorln("redirect to file must be the final operation in the sequence")
				return
			}

			if len(tokenGroup.Tokens) != 1 {
				Errorln("redirect must specify a single file path target")
				return
			}
		}

		if tokenGroup.FlowControl == FLOW_CONTROL_PIPE || tokenGroup.FlowControl == FLOW_CONTROL_UNSPECIFIED {
			isCapture := true
			if i == len(tokenGroups)-1 {
				isCapture = false
			}

			tokens := tokenGroup.Tokens
			command, parentFlags, remaining := c.LocateCommand(tokens)
			if command == nil {
				Errorln(fmt.Sprintf("unknown command \"%s\"", remaining[0]))
				return
			}

			isHelp := slices.Contains(tokens, "--help")
			if isHelp {
				fmt.Println(command.GetHelpString())
				return
			}

			argMap, err := command.ClassifyTokens(remaining, parentFlags)
			if err != nil {
				Errorln(err.Error())
				return
			}

			for _, arg := range command.Arguments {
				if _, exists := argMap[arg.Name]; !exists {
					Errorln(fmt.Sprintf("missing argument \"%s\"", arg.Name))
					return
				}
			}

			execSequence = append(execSequence, &BoundExec{
				IsCapturingOutput: isCapture,
				Command:           command,
				ArgMap:            argMap,
			})

		}
	}

	capturedBytes := []byte{}
	for _, bindExec := range execSequence {
		// capture each command's stdout and pass to the input of the next

		if bindExec.IsCapturingOutput {
			formerStdout := os.Stdout
			pipeRead, pipeWrite, err := os.Pipe()
			if err != nil {
				Errorln(err.Error())
				return
			}

			os.Stdout = pipeWrite
			completionChan := make(chan *Capture, 1)

			go func() {
				capture := &Capture{}
				_, err := io.Copy(&capture.Buffer, pipeRead)
				if err != nil {
					capture.Error = err
				}
				completionChan <- capture
			}()

			err = bindExec.Command.OnExecute(bindExec.Command, bindExec.ArgMap, capturedBytes)
			pipeWrite.Close()

			capture := <-completionChan
			os.Stdout = formerStdout

			if err != nil {
				Errorln(err.Error())
				return
			}

			if capture.Error != nil {
				Errorln(err.Error())
				return
			}

			// We don't capture terminal codes
			capturedBytes = term.StripTerminalEscapeSequences(capture.Buffer.Bytes())
			continue
		}

		err := bindExec.Command.OnExecute(bindExec.Command, bindExec.ArgMap, capturedBytes)
		if err != nil {
			Errorln(err.Error())
			return
		}
	}

	finalTokenGroup := tokenGroups[len(tokenGroups)-1]
	if finalTokenGroup.FlowControl == FLOW_CONTROL_REDIRECT {
		target := finalTokenGroup.Tokens[0]
		err := os.WriteFile(target, capturedBytes, 0644)
		if err != nil {
			Errorln("unable to write to file ", target, ": ", err.Error())
			return
		}
	}
}

func (c *Commander) Run() error {
	return c.shell.ReadUntilTerm()
}
