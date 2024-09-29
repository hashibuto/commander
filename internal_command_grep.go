package commander

import (
	"fmt"
	"strings"
)

const (
	PatternArg     string = "pattern"
	InsensitiveArg string = "insensitive"
)

var GrepCommand = &Command{
	Name:        "grep",
	Description: "filter and pattern match input",
	Arguments: []*Argument{
		{
			Name:        PatternArg,
			Description: "search pattern",
			ArgType:     ArgTypeString,
		},
	},
	Flags: []*Flag{
		{
			Name:         InsensitiveArg,
			ShortName:    "i",
			Description:  "case insensitive matching",
			ArgType:      ArgTypeBool,
			DefaultValue: false,
		},
	},
	OnExecute: func(c *Command, args ArgMap, capturedInput []byte) error {
		pattern := args.GetString(PatternArg)
		insensitive := args.GetBool(InsensitiveArg)
		lowerPattern := strings.ToLower(pattern)

		if len(capturedInput) > 0 {
			ci := string(capturedInput)
			lines := strings.Split(ci, "\n")
			for _, line := range lines {
				if insensitive {
					if !strings.Contains(strings.ToLower(line), lowerPattern) {
						continue
					}
				} else if !strings.Contains(line, pattern) {
					continue
				}

				fmt.Println(line)
			}
		}

		return nil
	},
}
