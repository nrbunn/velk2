package commands

import (
	"fmt"
	"log"
	"velk2/pkg/interfaces"
)

type Move struct {
	Server interfaces.ServerInterface
}

func (c *Move) Action(playerInterface interfaces.Character, command string, commandOptions ...string) {
	player := playerInterface

	targetRoom := player.GetRoom().GetExit(command)
	if targetRoom == nil {
		player.SendToCharacter("You can't go that way.\r\n")
		return
	}

	player.GetRoom().RemoveCharacter(player)
	c.sendToRoom(player.GetRoom(), fmt.Sprintf("%s has left.\r\n", player.GetName()))
	player.SetRoom(targetRoom.TargetRoom())
	c.sendToRoom(targetRoom.TargetRoom(), fmt.Sprintf("%s has arrived.\r\n", player.GetName()))
	targetRoom.TargetRoom().AddCharacter(player)

	message := fmt.Sprintf("%s\r\n%s\r\n%s\r\n", targetRoom.TargetRoom().GetName(), targetRoom.TargetRoom().GetDescription(), targetRoom.TargetRoom().GetExits())
	for _, character := range targetRoom.TargetRoom().GetCharacters() {
		message += fmt.Sprintf("%s\r\n", character.GetName())
	}

	player.SendToCharacter(message)

}

func (c *Move) sendToRoom(room interfaces.Room, msg string) {
	for _, character := range room.GetCharacters() {
		character.SendToCharacter(msg)
	}
}

func (c *Move) sendToAllPlayers(msg string) {
	for _, player := range c.Server.GetPlayers() {
		err := player.SendToCharacter(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
