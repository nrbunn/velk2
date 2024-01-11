package velk

import (
	"fmt"
	"log"
)

func Kill(player *Player, _ string, target ...string) {
	if len(target) == 0 {
		player.SendToCharacter("Kill who?\r\n")
		return
	}

	targetName := target[0]
	var targetMob *Mob
	for _, mob := range player.Room.Mobs {
		log.Printf("mob: %s", mob.Name)
		if mob.Name == targetName {
			targetMob = mob
			break
		}
	}

	if targetMob == nil {
		player.SendToCharacter(fmt.Sprintf("You don't see a %s here.\r\n", targetName))
		return
	}

	player.SendToCharacter(fmt.Sprintf("You attack %s!\r\n", targetMob.Name))
	player.AddFightingTarget(targetMob)
}
