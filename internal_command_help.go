package commander

import (
	"fmt"
	"sort"
)

var HelpCommand = &Command{
	Name:        "help",
	Description: "display contextual command help",
	OnExecute: func(c *Command, args ArgMap, capturedInput []byte) error {
		fmt.Println("Command list:")

		commandList := []string{}
		for _, cmd := range c.Commander.commandMap {
			commandList = append(commandList, cmd.Name)
		}
		sort.Slice(commandList, func(i, j int) bool {
			return commandList[i] < commandList[j]
		})

		for _, cmdName := range commandList {
			cmd := c.Commander.commandMap[cmdName]
			fmt.Printf("  %s%s\n", PadRight(cmd.Name, COMMAND_PADDING), cmd.Description)
		}

		return nil
	},
}
