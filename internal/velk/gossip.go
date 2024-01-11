package velk

import (
	"fmt"
	"log"
)

func Gossip(playerInterface *Player, _ string, commandOptions ...string) {
	var playerMessage string
	if len(commandOptions) > 0 {
		playerMessage = commandOptions[0]
	}

	messageSuffix := fmt.Sprintf("\"%s&n\"\r\n", playerMessage)

	for _, targetPlayer := range Players {
		var message string
		if targetPlayer == playerInterface {
			message = fmt.Sprintf("You gossip, %s", messageSuffix)
		} else {
			message = fmt.Sprintf("%s gossips, %s", playerInterface.Name, messageSuffix)
		}
		err := targetPlayer.SendToCharacter(message)
		if err != nil {
			log.Println(err)
		}
	}
}
