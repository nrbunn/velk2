package velk

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const clearScreenCommand = "\033[2J\033[H"

const (
	ReditSave                    RState = 1
	ReditMainMenu                RState = 2
	ReditMainMenuSelection       RState = 3
	ReditNameMenu                RState = 4
	ReditNameInput               RState = 5
	ReditDescription             RState = 6
	ReditDescriptionInput        RState = 7
	ReditExits                   RState = 8
	ReditExitsInput              RState = 9
	ReditExitEdit                RState = 10
	ReditExitEditInput           RState = 11
	ReditExitEditRoomSelect      RState = 12
	ReditExitEditNewRoom         RState = 13
	ReditExitEditRemoveRoom      RState = 14
	ReditExitEditRoomSelectInput RState = 15
)

type RState int

type ReditData struct {
	TargetRoom        *Room
	EditRoom          *Room
	EditDirection     string
	DescriptionBuffer []string
	State             RState
}

func Redit(player *Player, _ string, commandOptions ...string) {
	if len(commandOptions[0]) != 0 {
		//find room
		return
	}
	player.State = PlayerStateEditing
	player.Olc.Mode = ReditMode
	player.Olc.Redit = newReditData(player)
}

func newReditData(player *Player) *ReditData {
	editRoom := &Room{}
	*editRoom = *player.Room
	return &ReditData{
		TargetRoom:        player.Room,
		EditRoom:          editRoom,
		EditDirection:     "",
		State:             ReditMainMenu,
		DescriptionBuffer: strings.Split(editRoom.Description, "\r\n"),
	}
}

func handleRedit(player *Player) {

	switch player.Olc.Redit.State {
	case ReditMainMenu:
		handleReditMainMenu(player)
	case ReditMainMenuSelection:
		handleReditMainMenuSelection(player)
	case ReditNameMenu:
		handleReditNameMenu(player)
	case ReditNameInput:
		handleReditNameInput(player)
	case ReditDescription:
		handleReditDescription(player)
	case ReditDescriptionInput:
		handleReditDescriptionInput(player)
	case ReditExits:
		handleReditExits(player)
	case ReditExitsInput:
		handleReditExitsInput(player)
	case ReditExitEdit:
		handleReditExitEdit(player)
	case ReditExitEditInput:
		handleReditExitEditInput(player)
	case ReditExitEditRoomSelect:
		handleReditExitEditRoomSelect(player)
	case ReditExitEditNewRoom:
		handleReditExitEditNewRoom(player)
	case ReditExitEditRemoveRoom:
		handleReditExitEditRemoveRoom(player)
	case ReditExitEditRoomSelectInput:
		handleReditExitEditRoomSelectInput(player)
	case ReditSave:
		handleReditSave(player)
	}
}

func handleReditSave(player *Player) {
	player.Olc.Redit.TargetRoom.Name = player.Olc.Redit.EditRoom.Name
	player.Olc.Redit.TargetRoom.Description = player.Olc.Redit.EditRoom.Description
	for pair := player.Olc.Redit.EditRoom.Exits.Oldest(); pair != nil; pair = pair.Next() {
		if pair.Value.ZoneId == -2 && pair.Value.VirtualId == -2 {
			room := player.Olc.Redit.TargetRoom.Zone.CreateRoom()
			targetroomdir := getOppositeDirection(pair.Key)
			room.Exits.Set(targetroomdir, Vnum{ZoneId: player.Olc.Redit.EditRoom.Zone.Id, VirtualId: player.Olc.Redit.EditRoom.Id})
			err := room.Save()
			if err != nil {
				return
			}
			player.Olc.Redit.TargetRoom.Exits.Set(pair.Key, Vnum{ZoneId: room.Zone.Id, VirtualId: room.Id})
		} else {
			player.Olc.Redit.TargetRoom.Exits.Set(pair.Key, pair.Value)
		}

	}
	err := player.Olc.Redit.TargetRoom.Save()
	if err != nil {
		log.Printf("Error saving room: %v\n", err)
		player.SendToCharacter("Failed to save room")
		return
	}
	player.State = PlayerStateActive
	player.Olc.Mode = Inactive
	player.Olc.Redit = nil
	player.SendToCharacter("Room saved.\r\n")
}
func handleReditMainMenu(player *Player) {
	player.SendToCharacter(clearScreenCommand)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("(1)Name: %s\r\n", player.Olc.Redit.EditRoom.Name))
	builder.WriteString(fmt.Sprintf("(2)Description:\r\n%s\r\n", player.Olc.Redit.EditRoom.Description))
	builder.WriteString("(3)Exits:\r\n")
	for pair := player.Olc.Redit.EditRoom.Exits.Oldest(); pair != nil; pair = pair.Next() {
		if pair.Value.ZoneId != -1 && pair.Value.VirtualId != -1 {
			builder.WriteString(fmt.Sprintf("%s: %s\r\n", pair.Key, pair.Value.ToString()))
		} else {
			builder.WriteString(fmt.Sprintf("%s: None\r\n", pair.Key))
		}
	}
	builder.WriteString("(A)Abort\r\n")
	builder.WriteString("(S)Save\r\n")
	builder.WriteString("Selection: ")
	player.SendToCharacter(builder.String())
	player.Olc.Redit.State = ReditMainMenuSelection
}

func handleReditMainMenuSelection(player *Player) {
	option, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	option = strings.ToLower(strings.TrimSpace(option))
	switch option {
	case "1":
		player.Olc.Redit.State = ReditNameMenu
	case "2":
		player.Olc.Redit.State = ReditDescription
	case "3":
		player.Olc.Redit.State = ReditExits
	case "a":
		player.State = PlayerStateActive
		player.Olc.Mode = Inactive
		player.Olc.Redit = nil
	case "s":
		player.Olc.Redit.State = ReditSave
	}
}

func handleReditNameMenu(player *Player) {
	player.SendToCharacter(clearScreenCommand)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Name: %s\r\n", player.Olc.Redit.EditRoom.Name))
	builder.WriteString("Enter new name: ")
	player.SendToCharacter(builder.String())
	player.Olc.Redit.State = ReditNameInput
}

func handleReditNameInput(player *Player) {
	name, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)
	player.Olc.Redit.EditRoom.Name = name
	player.Olc.Redit.State = ReditMainMenu
}

func handleReditDescription(player *Player) {
	builder := strings.Builder{}
	builder.WriteString("Instructions: /s to save, /h for more options.\r\n")
	for _, line := range player.Olc.Redit.DescriptionBuffer {
		builder.WriteString(line)
		builder.WriteString("\r\n")
	}
	builder.WriteString("> ")
	player.SendToCharacter(builder.String())
	player.Olc.Redit.State = ReditDescriptionInput
}

func handleReditDescriptionInput(player *Player) {
	description, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	description = strings.TrimSpace(description)

	if len(description) == 2 {
		switch description {
		case "/s":
			player.Olc.Redit.EditRoom.Description = strings.Join(player.Olc.Redit.DescriptionBuffer, "\r\n")
			player.Olc.Redit.State = ReditMainMenu
			return
		case "/h":
			player.SendToCharacter("Instructions: /s to save, /h for more options.\r\n")
			return
		}
	}

	player.Olc.Redit.DescriptionBuffer = append(player.Olc.Redit.DescriptionBuffer, description)
}

func handleReditExits(player *Player) {
	builder := strings.Builder{}
	i := 1
	for pair := player.Olc.Redit.EditRoom.Exits.Oldest(); pair != nil; pair = pair.Next() {
		if pair.Value.VirtualId != -1 && pair.Value.ZoneId != -1 {
			builder.WriteString(fmt.Sprintf("(%d)%s: %s\r\n", i, pair.Key, pair.Value.ToString()))
		} else {
			builder.WriteString(fmt.Sprintf("(%d)%s: None\r\n", i, pair.Key))
		}
		i++
	}
	builder.WriteString("(A)Add Custom Exit\r\n")
	builder.WriteString("(B)Back\r\n")
	builder.WriteString("Selection: ")
	player.SendToCharacter(builder.String())
	player.Olc.Redit.State = ReditExitsInput
}

func handleReditExitsInput(player *Player) {
	option, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	option = strings.ToLower(strings.TrimSpace(option))
	switch option {
	case "a":
		player.Olc.Redit.State = ReditExitsInput
	case "b":
		player.Olc.Redit.State = ReditMainMenu
	default:
		optionNum, optionNumErr := strconv.Atoi(option)
		if optionNumErr != nil || optionNum > player.Olc.Redit.EditRoom.Exits.Len() {
			player.SendToCharacter("Invalid option.\r\n")
			player.Olc.Redit.State = ReditExits
			return
		}
		i := 1
		for pair := player.Olc.Redit.EditRoom.Exits.Oldest(); pair != nil; pair = pair.Next() {
			if i == optionNum {
				player.Olc.Redit.EditDirection = pair.Key
				break
			}
			i++
		}
		player.Olc.Redit.State = ReditExitEdit
	}
}

func handleReditExitEdit(player *Player) {
	var rnum Vnum
	var found bool
	if rnum, found = player.Olc.Redit.EditRoom.Exits.Get(player.Olc.Redit.EditDirection); !found {
		player.SendToCharacter("Exit not found.\r\n")
		player.Olc.Redit.State = ReditExits
		return
	}

	var roomId string
	if rnum.ZoneId == -1 && rnum.VirtualId == -1 {
		roomId = "None"
	} else if rnum.VirtualId == 0 {
		roomId = "New Room"
	} else {
		roomId = rnum.ToString()
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Exit: %s\r\n", player.Olc.Redit.EditDirection))
	builder.WriteString(fmt.Sprintf("(1)Room: %s\r\n", roomId))
	builder.WriteString("(N)New Room\r\n")
	builder.WriteString("(R)Remove Exit\r\n")
	builder.WriteString("(B)Back\r\n")
	builder.WriteString("Selection: ")
	player.SendToCharacter(builder.String())
	player.Olc.Redit.State = ReditExitEditInput
}

func handleReditExitEditInput(player *Player) {
	option, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	option = strings.ToLower(strings.TrimSpace(option))
	switch option {
	case "1":
		player.Olc.Redit.State = ReditExitEditRoomSelect
	case "n":
		player.Olc.Redit.State = ReditExitEditNewRoom
	case "r":
		player.Olc.Redit.State = ReditExitEditRemoveRoom
	case "b":
		player.Olc.Redit.State = ReditExits
	}
}

func handleReditExitEditRoomSelect(player *Player) {
	builder := strings.Builder{}
	builder.WriteString("Enter room id: ")
	player.SendToCharacter(builder.String())
	player.Olc.Redit.State = ReditExitEditRoomSelectInput
}

func handleReditExitEditRoomSelectInput(player *Player) {
	rnumStr, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}

	rnum, err := ParseVnum(rnumStr)
	if err != nil {
		player.SendToCharacter("Invalid room id.\r\n")
		player.Olc.Redit.State = ReditExitEdit
		return
	}
	if len(player.Olc.Redit.EditRoom.Zone.Rooms) < rnum.VirtualId {
		player.SendToCharacter("Invalid room id.\r\n")
		player.Olc.Redit.State = ReditExitEdit
		return
	}

	player.Olc.Redit.EditRoom.Exits.Set(player.Olc.Redit.EditDirection, rnum)
	player.Olc.Redit.State = ReditExits
}

func handleReditExitEditNewRoom(player *Player) {
	player.Olc.Redit.EditRoom.Exits.Set(player.Olc.Redit.EditDirection, Vnum{ZoneId: -2, VirtualId: -2})
	player.Olc.Redit.State = ReditExits
}

func handleReditExitEditRemoveRoom(player *Player) {
	player.Olc.Redit.EditRoom.Exits.Set(player.Olc.Redit.EditDirection, Vnum{ZoneId: -1, VirtualId: -1})
	player.Olc.Redit.State = ReditExits
}
