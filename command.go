package commander

import (
	"fmt"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Flags       []*Flag
	Arguments   []*Argument
	SubCommands []*Command

	commandMap map[string]*Command
	flagMap    map[string]*Flag
	argMap     map[string]*Argument
}

// Validate ensures that the command is valid, returning a descriptive error if it is not.
// Also, performs any run-time optimization.  This should only be called by the Commander.
func (c *Command) Validate() error {
	c.commandMap = map[string]*Command{}
	c.flagMap = map[string]*Flag{}
	c.argMap = map[string]*Argument{}

	for _, subCmd := range c.SubCommands {
		if _, exists := c.commandMap[subCmd.Name]; exists {
			return fmt.Errorf("sub-command \"%s\" under \"%s\" is defined multiple times", subCmd.Name, c.Name)
		}

		c.commandMap[subCmd.Name] = subCmd
	}

	for _, arg := range c.Arguments {
		if _, exists := c.flagMap[arg.Name]; exists {
			return fmt.Errorf("argument name \"%s\" on command \"%s\" is already defined as a flag", arg.Name, c.Name)
		}
		if _, exists := c.argMap[arg.Name]; exists {
			return fmt.Errorf("argument name \"%s\" on command \"%s\" is defined multiple times", arg.Name, c.Name)
		}
		err := arg.Validate()
		if err != nil {
			return err
		}
	}

	for _, flag := range c.Flags {
		if len(flag.Name) <= 1 {
			if _, exists := c.flagMap[flag.Name]; exists {
				return fmt.Errorf("flag name \"%s\" on command \"%s\" is defined multiple times", flag.Name, c.Name)
			}
		}

		if len(flag.ShortName) == 1 {
			if _, exists := c.flagMap[flag.ShortName]; exists {
				return fmt.Errorf("flag short name \"%s\" on command \"%s\" is defined multiple times", flag.ShortName, c.Name)
			}
		}
	}

	return nil
}

// ClassifyTokens attempts to classify the token array using the defined flags and arguments, in order to populate a name to value mapping
func (c *Command) ClassifyTokens(tokens []string) (map[string]any, error) {
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
				flagBody := t[1:]
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
				flag, ok := c.flagMap[name]
				if !ok {
					return nil, fmt.Errorf("unrecognized flag %s", name)
				}

				if hasValue {
					err := curFlag.PopulateMap(value, tokenMap)
					if err != nil {
						return nil, fmt.Errorf("invalid value for flag %s: %w", curFlag.GetInvocation(), err)
					}
				} else {
					if flag.ArgType != ArgTypeBool {
						curFlag = flag
					}
				}
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

	return tokenMap, nil
}
