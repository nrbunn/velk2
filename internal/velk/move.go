package velk

import (
	"fmt"
	"strings"
	"velk2/internal/utils"
)

func Move(player *Player, command string, _ ...string) {
	targetExit := getTargetExit(command)
	targetRnum := player.Room.GetExit(targetExit)

	if targetRnum.ZoneId == -1 {
		player.SendToCharacter("You can't go that way.\r\n")
		return
	}

	targetRoom := Zones[targetRnum.ZoneId].GetRoom(targetRnum.VirtualId)
	player.Room.RemoveCharacter(player)
	sendToRoom(player.Room, fmt.Sprintf("%s has left.\r\n", player.Name))
	player.Room = targetRoom
	sendToRoom(targetRoom, fmt.Sprintf("%s has arrived.\r\n", player.Name))
	targetRoom.AddCharacter(player)

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
