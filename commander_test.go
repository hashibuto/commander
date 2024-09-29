package commander

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CommanderTestSuite struct {
	suite.Suite

	TheCommander *Commander
}

type AnimalType string

const (
	AnimalTypeMammal  AnimalType = "mammal"
	AnimalTypeReptile AnimalType = "reptile"
	AnimalTypeBird    AnimalType = "bird"
)

func (suite *CommanderTestSuite) SetupTest() {
	var err error
	suite.TheCommander, err = NewCommander(Config{
		Commands: []*Command{
			{
				Name:        "farm",
				Description: "interact with the farm",
				SubCommands: []*Command{
					{
						Name:        "snapshot",
						Description: "work with a snapshot",
						Flags: []*Flag{
							{
								Name:        "type",
								Description: "snapshot type",
								ShortName:   "t",
								ArgType:     ArgTypeString,
								OneOf:       []any{"image", "inventory"},
							},
						},
						SubCommands: []*Command{
							{
								Name:        "create",
								Description: "create a snapshot",
								Flags: []*Flag{
									{
										Name:         "timestamp",
										ShortName:    "s",
										Description:  "set timestamp",
										ArgType:      ArgTypeBool,
										DefaultValue: true,
									},
								},
								OnExecute: func(c *Command, args ArgMap, capturedInput []byte) error {
									return nil
								},
							},
							{
								Name:        "update",
								Description: "update a snapshot",
								OnExecute: func(c *Command, args ArgMap, capturedInput []byte) error {
									return nil
								},
							},
						},
					},
					{
						Name:        "inventory",
						Description: "obtain animal inventory",
						Flags: []*Flag{
							{
								Name:      "sort",
								ShortName: "s",
							},
						},
						OnExecute: func(c *Command, args ArgMap, capturedInput []byte) error {
							return nil
						},
					},
					{
						Name:        "add",
						Description: "add an animal to the farm",
						Flags: []*Flag{
							{
								Name:       "type",
								ShortName:  "t",
								ArgType:    ArgTypeString,
								OneOf:      []any{AnimalTypeMammal, AnimalTypeBird, AnimalTypeReptile},
								IsRequired: true,
							},
						},
						OnExecute: func(c *Command, args ArgMap, capturedInput []byte) error {
							return nil
						},
					},
				},
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (suite *CommanderTestSuite) TestLocateCommand() {
	tokens := []string{"farm", "add", "-t", "mammal"}
	command, _, commandTokens := suite.TheCommander.LocateCommand(tokens)
	assert.NotNil(suite.T(), command)
	assert.Len(suite.T(), commandTokens, 2)
}

func (suite *CommanderTestSuite) TestLocateDeepCommand() {
	tokens := []string{"farm", "snapshot", "create"}
	command, _, _ := suite.TheCommander.LocateCommand(tokens)
	assert.NotNil(suite.T(), command)
}

func TestCommander(t *testing.T) {
	suite.Run(t, new(CommanderTestSuite))
}
