package structs

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"strings"
	"velk2/pkg/interfaces"
)

type Commands struct {
	CommandMap *orderedmap.OrderedMap[string, interfaces.CommandInterface]
}

func NewCommands() *Commands {
	return &Commands{CommandMap: orderedmap.New[string, interfaces.CommandInterface]()}
}
func (c *Commands) addCommand(commandString string, command interfaces.CommandInterface) {
	c.CommandMap.Set(commandString, command)
}

func (c *Commands) getCommand(commandString string) interfaces.CommandInterface {

	for pair := c.CommandMap.Oldest(); pair != nil; pair = pair.Next() {
		if strings.HasPrefix(pair.Key, commandString) {
			return pair.Value
		}
	}

	return nil
}
