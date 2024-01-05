package commands

import (
	"fmt"
	"log"
	"velk2/pkg/interfaces"
)

type Gossip struct {
	Server interfaces.ServerInterface
}

func (c *Gossip) Action(playerInterface interfaces.Character, command string, commandOptions ...string) {
	player := playerInterface
	message := ""
	if len(commandOptions) > 0 {
		message = commandOptions[0]
	}
	sendMessage := fmt.Sprintf("%s(%s) gossips, \"%s&n\"\r\n", player.GetName(), player.GetUUID(), message)

	c.sendToAllPlayers(sendMessage)
}

func (c *Gossip) sendToAllPlayers(msg string) {
	for _, player := range c.Server.GetPlayers() {
		err := player.SendToCharacter(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
