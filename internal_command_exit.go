package commander

import (
	ns "github.com/hashibuto/nilshell"
)

var ExitCommand = &Command{
	Name:        "exit",
	Description: "exit the shell",
	OnExecute: func(c *Command, args map[string]any, capturedInput []byte) error {
		return ns.ErrEof
	},
}
