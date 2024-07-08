package commander

import (
	"fmt"
	"strings"

	ns "github.com/hashibuto/nilshell"
)

type Command struct {
	Name        string
	Description string
	Flags       []*Flag
	Arguments   []*Argument
	SubCommands []*Command
	OnExecute   func(c *Command, args map[string]any, capturedInput []byte) error

	Commander  *Commander
	commandMap map[string]*Command
	flagMap    map[string]*Flag
	argMap     map[string]*Argument
}

// Validate ensures that the command is valid, returning a descriptive error if it is not.
// Also, performs any run-time optimization.  This should only be called by the Commander.
func (c *Command) Validate(parentFlags map[string]struct{}) error {
	if parentFlags == nil {
		parentFlags = map[string]struct{}{}
	}

	c.commandMap = map[string]*Command{}
	c.flagMap = map[string]*Flag{}
	c.argMap = map[string]*Argument{}

	if len(c.SubCommands) > 0 && len(c.Arguments) > 0 {
		return fmt.Errorf("command \"%s\" cannot contain both subcommands and positional arguments", c.Name)
	}

	if len(c.SubCommands) > 0 && c.OnExecute != nil {
		return fmt.Errorf("command \"%s\" cannot contain both subcommands and an OnExecute handler", c.Name)
	}

	if len(c.SubCommands) == 0 && c.OnExecute == nil {
		return fmt.Errorf("command \"%s\" does not implement an OnExecute handler", c.Name)
	}

	for _, arg := range c.Arguments {
		if _, exists := parentFlags[arg.Name]; exists {
			return fmt.Errorf("argument name \"%s\" on command \"%s\" is already defined as parent command flag", arg.Name, c.Name)
		}

		if _, exists := c.flagMap[arg.Name]; exists {
			return fmt.Errorf("argument name \"%s\" on command \"%s\" is already defined as a flag", arg.Name, c.Name)
		}
		if _, exists := c.argMap[arg.Name]; exists {
			return fmt.Errorf("argument name \"%s\" on command \"%s\" is defined multiple times", arg.Name, c.Name)
		}
		err := arg.Validate()
		if err != nil {
			return fmt.Errorf("command \"%s\" - %w", c.Name, err)
		}

		c.argMap[arg.Name] = arg
	}

	helpFlag := &Flag{
		Name:        "help",
		ArgType:     ArgTypeBool,
		Description: "display contextual help",
	}

	// this is required b/c of singly defined internal commands which get "validated" more than once in the unit tests, and thus get multiple
	// --help flags otherwise
	found := false
	for _, flag := range c.Flags {
		if flag.Name == helpFlag.Name && flag.Description == helpFlag.Description {
			found = true
			break
		}
	}
	if !found {
		c.Flags = append(c.Flags, helpFlag)
	}

	for _, flag := range c.Flags {
		err := flag.Validate()
		if err != nil {
			return fmt.Errorf("command \"%s\" - %w", c.Name, err)
		}

		if flag.Name != "" {
			if _, exists := parentFlags[flag.Name]; exists && flag.Name != "help" {
				return fmt.Errorf("flag name \"%s\" on command \"%s\" is already defined as parent command flag", flag.Name, c.Name)
			}

			if _, exists := c.flagMap[flag.Name]; exists {
				return fmt.Errorf("flag name \"%s\" on command \"%s\" is defined multiple times", flag.Name, c.Name)
			}

			c.flagMap[flag.Name] = flag
			parentFlags[flag.Name] = struct{}{}
		}

		if flag.ShortName != "" {
			if _, exists := parentFlags[flag.ShortName]; exists {
				return fmt.Errorf("flag short name \"%s\" on command \"%s\" is already defined as parent command flag %+v", flag.ShortName, c.Name, parentFlags)
			}

			if _, exists := c.flagMap[flag.ShortName]; exists {
				return fmt.Errorf("flag short name \"%s\" on command \"%s\" is defined multiple times", flag.ShortName, c.Name)
			}

			c.flagMap[flag.ShortName] = flag
			parentFlags[flag.ShortName] = struct{}{}
		}
	}

	for _, subCmd := range c.SubCommands {
		if _, exists := c.commandMap[subCmd.Name]; exists {
			return fmt.Errorf("sub-command \"%s\" under \"%s\" is defined multiple times", subCmd.Name, c.Name)
		}

		c.commandMap[subCmd.Name] = subCmd
		parentCopy := map[string]struct{}{}
		for k, v := range parentFlags {
			parentCopy[k] = v
		}
		err := subCmd.Validate(parentCopy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Command) Suggest(tokens []string, parentFlags []*Flag) []*ns.AutoComplete {
	argNum := 0

	allFlags := append(parentFlags, c.Flags...)

	noFlags := false
	var curFlag *Flag
	for idx, t := range tokens {
		isFinal := idx == len(tokens)-1
		if curFlag != nil {
			if isFinal {
				// provide suggestions for this flag's value set if any
				return curFlag.SuggestValues(t)
			}

			curFlag = nil
			continue
		}

		if !noFlags {
			if t == "--" && idx < len(tokens)-1 {
				noFlags = true
				continue
			}

			if strings.HasPrefix(t, "-") && !strings.HasPrefix(t, "--") {
				flagBody := t[1:]
				// if value is already being assigned
				if strings.Contains(flagBody, "=") {
					continue
				}
				prefix := flagBody

				if isFinal {
					sugg := []*ns.AutoComplete{}
					for _, f := range allFlags {
						if strings.HasPrefix(f.ShortName, prefix) {
							sugg = append(sugg, &ns.AutoComplete{
								Value:   fmt.Sprintf("-%s", f.ShortName),
								Display: fmt.Sprintf("%s  %s", f.GetInvocation(), f.Description),
							})
						}
					}

					return sugg
				}

				if f, ok := c.flagMap[flagBody]; ok {
					curFlag = f
				}
				continue
			}

			if strings.HasPrefix(t, "--") {
				flagBody := t[2:]
				// if value is already being assigned
				if strings.Contains(flagBody, "=") {
					continue
				}
				prefix := flagBody
				if isFinal {
					sugg := []*ns.AutoComplete{}
					for _, f := range allFlags {
						if strings.HasPrefix(f.Name, prefix) {
							sugg = append(sugg, &ns.AutoComplete{
								Value:   fmt.Sprintf("--%s", f.Name),
								Display: fmt.Sprintf("%s  %s", f.GetInvocation(), f.Description),
							})
						}
					}

					return sugg
				}

				if f, ok := c.flagMap[flagBody]; ok {
					curFlag = f
				}
				continue
			}
		}

		if len(c.SubCommands) > 0 {
			sugg := []*ns.AutoComplete{}
			for _, sub := range c.SubCommands {
				if strings.HasPrefix(sub.Name, t) {
					sugg = append(sugg, &ns.AutoComplete{
						Value:   sub.Name,
						Display: sub.Name,
					})
				}
			}
			return sugg
		}

		if len(c.Arguments) == 0 {
			return nil
		}

		var curArg *Argument
		if argNum >= len(c.Arguments) {
			if !c.Arguments[len(c.Arguments)-1].AllowMultiple {
				return nil
			}

			curArg = c.Arguments[len(c.Arguments)-1]
		} else {
			curArg = c.Arguments[argNum]
		}

		if isFinal {
			return curArg.SuggestValues(t)
		}
		argNum++
	}

	return nil
}

// ClassifyTokens attempts to classify the token array using the defined flags and arguments, in order to populate a name to value mapping
func (c *Command) ClassifyTokens(tokens []string, parentFlags []*Flag) (map[string]any, error) {
	allFlagMap := map[string]*Flag{}
	for k, v := range c.flagMap {
		allFlagMap[k] = v
	}
	for _, p := range parentFlags {
		if p.Name != "" {
			allFlagMap[p.Name] = p
		}
		if p.ShortName != "" {
			allFlagMap[p.ShortName] = p
		}
	}

	tokenMap := map[string]any{}
	argNum := 0

	noFlags := false
	var curFlag *Flag
	for _, t := range tokens {
		name := ""
		value := ""
		hasValue := false

		if curFlag != nil {
			// Grab the value for the active flag
			err := curFlag.PopulateMap(t, tokenMap)
			if err != nil {
				return nil, fmt.Errorf("invalid value for flag %s: %w", curFlag.GetInvocation(), err)
			}

			curFlag = nil
			continue
		}

		if !noFlags {
			if t == "--" {
				noFlags = true
				continue
			}

			if strings.HasPrefix(t, "-") && !strings.HasPrefix(t, "--") {
				flagBody := t[1:]
				parts := strings.SplitN(flagBody, "=", 2)
				name = parts[0]
				if len(parts) == 2 {
					hasValue = true
					value = parts[1]
				} else {
					hasValue = false
				}

				if len(name) > 1 {
					return nil, fmt.Errorf("malformed flag %s, did you mean -%s", t, t)
				}

				if len(name) == 0 {
					return nil, fmt.Errorf("missing flag name")
				}
			}

			if strings.HasPrefix(t, "--") && len(t) > 2 {
				flagBody := t[2:]
				parts := strings.SplitN(flagBody, "=", 2)
				name = parts[0]
				if len(parts) == 2 {
					hasValue = true
					value = parts[1]
				} else {
					hasValue = false
				}

				if len(name) == 1 {
					return nil, fmt.Errorf("malformed flag %s, did you mean -%s", t, name)
				}
			}

			// Flag detected
			if len(name) > 0 {
				flag, ok := allFlagMap[name]
				if !ok {
					return nil, fmt.Errorf("unrecognized flag %s", name)
				}

				if hasValue {
					err := curFlag.PopulateMap(value, tokenMap)
					if err != nil {
						return nil, fmt.Errorf("invalid value for flag %s: %w", curFlag.GetInvocation(), err)
					}
					continue
				}

				if flag.ArgType != ArgTypeBool {
					curFlag = flag
					continue
				}

				var v string
				if flag.DefaultValue == false {
					v = "true"
				} else {
					v = "false"
				}
				err := flag.PopulateMap(v, tokenMap)
				if err != nil {
					return nil, fmt.Errorf("invalid value for flag %s: %w", flag.GetInvocation(), err)
				}

				continue
			}
		}

		if len(c.Arguments) == 0 {
			return nil, fmt.Errorf("command \"%s\" does not accept any positional arguments", c.Name)
		}

		var curArg *Argument
		if argNum >= len(c.Arguments) {
			if !c.Arguments[len(c.Arguments)-1].AllowMultiple {
				return nil, fmt.Errorf("too many positional arguments provided")
			}

			curArg = c.Arguments[len(c.Arguments)-1]
		} else {
			curArg = c.Arguments[argNum]
		}

		err := curArg.PopulateMap(t, tokenMap)
		if err != nil {
			return nil, fmt.Errorf("invalid value for argument %s: %w", curArg.Name, err)
		}
		argNum++
	}

	// Apply all other tokens to the map
	for _, flag := range allFlagMap {
		err := flag.PopulateDefault(tokenMap)
		if err != nil {
			return nil, fmt.Errorf("command \"%s\" - %s", c.Name, err.Error())
		}
	}

	return tokenMap, nil
}

func (c *Command) getInvocation() string {
	parts := []string{
		c.Name,
	}

	if len(c.SubCommands) > 0 {
		parts = append(parts, "<subcommand>")
	}

	if len(c.Flags) > 0 {
		parts = append(parts, "[flags...]")
	}

	for _, arg := range c.Arguments {
		var text string
		if arg.AllowMultiple {
			text = fmt.Sprintf("<%s...>", arg.Name)
		} else {
			text = fmt.Sprintf("<%s>", arg.Name)
		}

		parts = append(parts, text)
	}

	return strings.Join(parts, " ")
}

func (c *Command) GetHelpString(parentFlags []*Flag) string {
	lines := []string{"Invocation:", c.getInvocation()}

	if len(c.SubCommands) > 0 {
		lines = append(lines, "", "Subcommands:")
		for _, sub := range c.SubCommands {
			lines = append(lines, fmt.Sprintf("  %s%s", PadRight(sub.Name, COMMAND_PADDING), sub.Description))
		}
	}

	if len(c.Arguments) > 0 {
		lines = append(lines, "", "Arguments:")
		for _, arg := range c.Arguments {
			description := []string{}
			if len(arg.Description) > 0 {
				description = append(description, arg.Description)
			}
			if arg.OneOf != nil {
				oneOf := []string{}
				for _, one := range arg.OneOf {
					oneOf = append(oneOf, fmt.Sprintf("%s", one))
				}
				description = append(description, fmt.Sprintf("one of %s", strings.Join(oneOf, ", ")))
			}
			lines = append(lines, fmt.Sprintf("  %s%s", PadRight(arg.Name, COMMAND_PADDING), strings.Join(description, " - ")))
		}
	}

	for i, flags := range [][]*Flag{parentFlags, c.Flags} {
		if len(flags) == 0 {
			continue
		}

		var label string
		if i == 0 {
			label = "Inherited flags:"
		} else {
			label = "Flags:"
		}
		lines = append(lines, "", label)
		for _, flag := range flags {
			description := []string{}
			if len(flag.Description) > 0 {
				description = append(description, flag.Description)
			}
			if flag.OneOf != nil {
				oneOf := []string{}
				for _, one := range flag.OneOf {
					oneOf = append(oneOf, fmt.Sprintf("\"%s\"", one))
				}
				description = append(description, fmt.Sprintf("one of %s", strings.Join(oneOf, ", ")))
				if flag.DefaultValue != nil {
					description = append(description, fmt.Sprintf("defaults to \"%s\"", flag.DefaultValue))
				}
			}
			lines = append(lines, fmt.Sprintf("  %s%s", PadRight(flag.GetPaddedInvocation(), COMMAND_PADDING), strings.Join(description, " - ")))
		}
	}

	return strings.Join(lines, "\n")
}
