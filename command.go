package commander

import (
	"fmt"
	"strings"
)

type Command struct {
	Name        string
	Description string
	Flags       []*Flag
	Argument    *Argument
	SubCommands []*Command

	commandMap map[string]*Command
	flagMap    map[string]*Flag
}

// Validate ensures that the command is valid, returning a descriptive error if it is not.
// Also, performs any run-time optimization.  This should only be called by the Commander.
func (c *Command) Validate() error {
	c.commandMap = map[string]*Command{}
	c.flagMap = map[string]*Flag{}

	for _, subCmd := range c.SubCommands {
		if _, exists := c.commandMap[subCmd.Name]; exists {
			return fmt.Errorf("sub-command \"%s\" under \"%s\" is defined multiple times", subCmd.Name, c.Name)
		}

		c.commandMap[subCmd.Name] = subCmd
	}

	if c.Argument != nil {
		err := c.Argument.Validate()
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

func (c *Command) ClassifyTokens(tokens []string) (map[string]any, []any, error) {
	args := []any{}
	flagMap := map[string]any{}

	noFlags := false
	var curFlag *Flag
	for _, t := range tokens {
		name := ""
		value := ""
		hasValue := false
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
					return nil, nil, fmt.Errorf("malformed flag %s, did you mean -%s", t, t)
				}

				if len(name) == 0 {
					return nil, nil, fmt.Errorf("missing flag name")
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
					return nil, nil, fmt.Errorf("malformed flag %s, did you mean -%s", t, name)
				}
			}

			// Flag detected
			if len(name) > 0 {
				flag, ok := c.flagMap[name]
				if !ok {
					return nil, nil, fmt.Errorf("unrecognized flag %s", name)
				}

				if hasValue {
					if flag.ShortName != "" {
						flagMap[flag.ShortName] = value
					}

					if flag.Name != "" {
						flagMap[flag.Name] = value
					}
				} else {
					curFlag = flag
				}
			}
		}
	}

	return nil, nil
}
