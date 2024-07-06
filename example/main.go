package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashibuto/commander"
	"github.com/hashibuto/nilshell/pkg/term"
	"gopkg.in/yaml.v3"
)

type ResourceType string

const (
	Process      ResourceType = "process"
	ProcessGroup ResourceType = "group"
)

const (
	ResourceTypeArg string = "resource-type"
	OutputFlag      string = "output"
)

func getProcesses(outputType string) error {
	if outputType == "table" {
		commander.Println(
			term.PadRight(commander.Sprintf(commander.C_BOLD, "PID"), 12, 2),
			term.PadRight(commander.Sprintf(commander.C_BOLD, "NAME"), 25, 2),
			commander.Sprintf(commander.C_BOLD, "INVOCATION"),
		)
		for _, procObj := range ProcessList {
			commander.Println(
				term.PadRight(fmt.Sprintf("%d", procObj.Id), 12, 2),
				term.PadRight(procObj.Name, 25, 2),
				procObj.Invocation,
			)
		}
	} else if outputType == "json" {
		jBytes, _ := json.Marshal(ProcessList)
		fmt.Println(string(jBytes))
	} else {
		// must be yaml
		yBytes, _ := yaml.Marshal(ProcessList)
		fmt.Println(string(yBytes))
	}

	return nil
}

func getProcessGroups(outputType string) error {
	if outputType == "table" {
		commander.Println(
			term.PadRight(commander.Sprintf(commander.C_BOLD, "GROUP"), 20, 2),
			term.PadRight(commander.Sprintf(commander.C_BOLD, "PID"), 12, 2),
			commander.Sprintf(commander.C_BOLD, "PROCESS"),
		)
		for _, groupObj := range ProcessGroups {
			for _, processObj := range groupObj.Processes {
				commander.Println(
					term.PadRight(groupObj.Name, 20, 2),
					term.PadRight(fmt.Sprintf("%d", processObj.Id), 12, 2),
					processObj.Name,
				)
			}
		}
	} else if outputType == "json" {
		jBytes, _ := json.Marshal(ProcessGroups)
		fmt.Println(string(jBytes))
	} else {
		// must be yaml
		yBytes, _ := yaml.Marshal(ProcessGroups)
		fmt.Println(string(yBytes))
	}

	return nil
}

func main() {
	c, err := commander.NewCommander(commander.Config{
		Prompt: commander.Sprintf(commander.FgColor(168, 94, 29), "demo", commander.FgColor(255, 235, 15), " » "),
		Commands: []*commander.Command{
			{
				Name:        "get",
				Description: "get information about a resource",
				Arguments: []*commander.Argument{
					{
						Name:        ResourceTypeArg,
						Description: "type of resource to locate",
						OneOf:       []any{Process, ProcessGroup},
					},
				},
				Flags: []*commander.Flag{
					{
						Name:         OutputFlag,
						Description:  "command output format",
						ArgType:      commander.ArgTypeString,
						ShortName:    "o",
						DefaultValue: "table",
						OneOf:        []any{"json", "yaml", "table"},
					},
				},
				OnExecute: func(c *commander.Command, args map[string]any, capturedInput []byte) error {
					resourceType := ResourceType(args[ResourceTypeArg].(string))
					outputType := args[OutputFlag].(string)

					var err error
					if resourceType == Process {
						err = getProcesses(outputType)
					} else {
						err = getProcessGroups(outputType)
					}

					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("exiting...")
}
