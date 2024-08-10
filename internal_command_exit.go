package commander

import "os"

var ExitCommand = &Command{
	Name:        "exit",
	Description: "exit the shell",
	OnExecute: func(c *Command, args map[string]any, capturedInput []byte) error {
		c.Commander.shell.Shutdown()
		os.Exit(0)

		return nil
	},
}
