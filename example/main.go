package main

import (
	"log"

	"github.com/hashibuto/commander"
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

func main() {
	c, err := commander.NewCommander(commander.Config{
		Commands: []*commander.Command{
			{
				Name: "get",
				Arguments: []*commander.Argument{
					{
						Name:        ResourceTypeArg,
						Description: "type of resource to locate",
						OneOf:       []any{Process, ProcessGroup},
					},
				},
				Flags: []*commander.Flag{
					{
						Name:      OutputFlag,
						ShortName: "o",
						OneOf:     []any{"json", "yaml"},
					},
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
}
