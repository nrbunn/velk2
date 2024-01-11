package velk

import (
	"fmt"
	"strings"
)

func Look(player *Player, _ string, _ ...string) {
	room := player.Room
	if room == nil {
		return
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("&c%s&n\r\n", room.Name))
	builder.WriteString(fmt.Sprintf("&w%s&n\r\n", room.Description))
	builder.WriteString(fmt.Sprintf("&c%s&n\r\n", printExits(room)))
	for _, character := range room.Characters {
		if character == player {
			continue
		}
		builder.WriteString(fmt.Sprintf("&y%s&n\r\n", character.Name))
	}
	for _, mob := range room.Mobs {
		builder.WriteString(fmt.Sprintf("&y%s&n\r\n", mob.ShortDesc))
	}
	player.SendToCharacter(builder.String())
}

func printExits(room *Room) (exitsStr string) {
	exitsStr = "[ "
	for pair := room.Exits.Oldest(); pair != nil; pair = pair.Next() {
		if pair.Value != nil {
			exitsStr += pair.Key + " "
		}
	}
	exitsStr += "]"
	return
}
