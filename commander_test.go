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
						Name:        "inventory",
						Description: "obtain animal inventory",
						Flags: []*Flag{
							{
								Name:      "sort",
								ShortName: "s",
							},
						},
					},
					{
						Name:        "add",
						Description: "add an animal to the farm",
						Flags: []*Flag{
							{
								Name:       "type",
								ShortName:  "t",
								OneOf:      []any{AnimalTypeMammal, AnimalTypeBird, AnimalTypeReptile},
								IsRequired: true,
							},
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

func TestCommander(t *testing.T) {
	suite.Run(t, new(CommanderTestSuite))
}
