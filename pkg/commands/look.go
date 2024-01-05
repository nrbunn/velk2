package commands

import (
	"fmt"
	"velk2/pkg/interfaces"
)

type Look struct {
	Server interfaces.ServerInterface
}

func (c *Look) Action(playerInterface interfaces.Character, command string, commandOptions ...string) {
	player := playerInterface
	room := player.GetRoom()
	if room == nil {
		return
	}

	message := fmt.Sprintf("%s\r\n%s\r\n%s\r\n", room.GetName(), room.GetDescription(), room.GetExits())
	for _, character := range room.GetCharacters() {
		message += fmt.Sprintf("%s\r\n", character.GetName())
	}

	player.SendToCharacter(message)
}
