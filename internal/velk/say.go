package velk

import (
	"fmt"
	"log"
)

func Say(player *Player, _ string, commandOptions ...string) {
	var playerMessage string
	if len(commandOptions) > 0 {
		playerMessage = commandOptions[0]
	}

	messageSuffix := fmt.Sprintf("\"%s&n\"\r\n", playerMessage)

	for _, targetPlayer := range player.Room.Characters {
		var message string
		if targetPlayer == player {
			message = fmt.Sprintf("You say, %s", messageSuffix)
		} else {
			message = fmt.Sprintf("%s says, %s", player.Name, messageSuffix)
		}
		err := targetPlayer.SendToCharacter(message)
		if err != nil {
			log.Println(err)
		}
	}
}
