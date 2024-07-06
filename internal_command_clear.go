package commander

var ClearCommand = &Command{
	Name:        "clear",
	Description: "clear the terminal",
	OnExecute: func(c *Command, args map[string]any, capturedInput []byte) error {
		c.Commander.shell.Clear()

		return nil
	},
}
