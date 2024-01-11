package velk

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"strings"
)

type CommandFunc func(player *Player, command string, commandOptions ...string)

type Commands struct {
	CommandMap *orderedmap.OrderedMap[string, CommandFunc]
}

func NewCommands() *Commands {
	return &Commands{CommandMap: orderedmap.New[string, CommandFunc]()}
}
func (c *Commands) AddCommand(commandString string, command CommandFunc) {
	c.CommandMap.Set(commandString, command)
}

func (c *Commands) GetCommand(commandString string) CommandFunc {

	for pair := c.CommandMap.Oldest(); pair != nil; pair = pair.Next() {
		if strings.HasPrefix(pair.Key, commandString) {
			return pair.Value
		}
	}

	return nil
}
