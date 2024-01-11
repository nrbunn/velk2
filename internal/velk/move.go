package velk

import (
	"fmt"
	"strings"
	"velk2/internal/utils"
)

func Move(player *Player, command string, _ ...string) {
	targetExit := getTargetExit(command)
	targetRoom := player.Room.GetExit(targetExit)

	if targetRoom == nil {
		player.SendToCharacter("You can't go that way.\r\n")
		return
	}

	player.Room.RemoveCharacter(player)
	sendToRoom(player.Room, fmt.Sprintf("%s has left.\r\n", player.Name))
	player.Room = targetRoom.TargetRoom()
	sendToRoom(targetRoom.TargetRoom(), fmt.Sprintf("%s has arrived.\r\n", player.Name))
	targetRoom.TargetRoom().AddCharacter(player)

	Look(player, "")

}

func getTargetExit(command string) string {
	for _, direction := range utils.Directions {
		if strings.HasPrefix(direction, command) {
			return direction
		}
	}
	return command
}

func sendToRoom(room *Room, msg string) {
	for _, character := range room.Characters {
		character.SendToCharacter(msg)
	}
}
