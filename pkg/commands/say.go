package commands

import (
	"fmt"
	"log"
	"velk2/pkg/interfaces"
)

type Say struct {
	Server interfaces.ServerInterface
}

func (c *Say) Action(playerInterface interfaces.Character, command string, commandOptions ...string) {
	player := playerInterface
	message := ""
	if len(commandOptions) > 0 {
		message = commandOptions[0]
	}
	sendMessage := fmt.Sprintf("%s(%s) says, \"%s&n\"\r\n", player.GetName(), player.GetUUID(), message)

	c.sendToAllPlayersInRoom(sendMessage, player.GetRoom())
}

func (c *Say) sendToAllPlayersInRoom(msg string, room interfaces.Room) {
	for _, player := range room.GetCharacters() {
		err := player.SendToCharacter(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
