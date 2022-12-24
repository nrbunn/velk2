package commands

import (
	"fmt"
	"log"
	"velk2/pkg/interfaces"
)

type Say struct {
	Server interfaces.ServerInterface
}

func (say *Say) Action(playerInterface interfaces.CharacterInterface, command string, commandOptions ...string) {
	player := playerInterface
	message := ""
	if len(commandOptions) > 0 {
		message = commandOptions[0]
	}
	sendMessage := fmt.Sprintf("%s says, \"%s&n\"\r\n", player.GetUUID(), message)

	say.sendToAllPlayers(sendMessage)
}

func (say *Say) sendToAllPlayers(msg string) {
	for _, player := range say.Server.GetPlayers() {
		err := player.SendToCharacter(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
