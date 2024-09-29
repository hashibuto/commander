package commander

import "github.com/hashibuto/nilshell/pkg/termutils"

var ClearCommand = &Command{
	Name:        "clear",
	Description: "clear the terminal",
	OnExecute: func(c *Command, args map[string]any, capturedInput []byte) error {
		termutils.ClearTerminal()
		termutils.SetCursorPos(1, 1)
		return nil
	},
}
