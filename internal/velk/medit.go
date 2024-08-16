package velk

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	MeditSave              MState = 1
	MeditMainMenu          MState = 2
	MeditMainMenuSelection MState = 3
	MeditNameMenu          MState = 4
	MeditNameInput         MState = 5
	MeditDescription       MState = 6
	MeditDescriptionInput  MState = 7
	MeditHealth            MState = 8
	MeditHealthInput       MState = 9
	MeditMana              MState = 10
	MeditManaInput         MState = 11
	MeditMove              MState = 12
	MeditMoveInput         MState = 13
)

type MState int

type MeditData struct {
	EditMob           *Mob
	DescriptionBuffer []string
	State             MState
}

func Medit(player *Player, _ string, commandOptions ...string) {
	if len(commandOptions[0]) == 0 {
		player.SendToCharacter("medit <vnum> or medit new")
		return
	}

	editMob := &Mob{}
	if commandOptions[0] == "new" {
		mob := NewMob(-1, player.Room.Zone)
		*editMob = *mob
	} else if vnum, err := ParseVnum(commandOptions[0]); err == nil {
		mob := Zones[vnum.ZoneId].Mobs[vnum.VirtualId]
		*editMob = *mob
	} else {
		player.SendToCharacter("medit <vnum> or medit new")
		return
	}
	player.State = PlayerStateEditing
	player.Olc.Mode = MeditMode
	player.Olc.Medit = newMeditData(player, editMob)
}

func newMeditData(player *Player, editMob *Mob) *MeditData {

	return &MeditData{
		EditMob:           editMob,
		State:             MeditMainMenu,
		DescriptionBuffer: strings.Split(editMob.Description, "\r\n"),
	}
}

func handleMedit(player *Player) {

	switch player.Olc.Medit.State {
	case MeditMainMenu:
		handleMeditMainMenu(player)
	case MeditMainMenuSelection:
		handleMeditMainMenuSelection(player)
	case MeditNameMenu:
		handleMeditNameMenu(player)
	case MeditNameInput:
		handleMeditNameInput(player)
	case MeditDescription:
		handleMeditDescription(player)
	case MeditDescriptionInput:
		handleMeditDescriptionInput(player)
	case MeditHealth:
		handleMeditHealthMenu(player)
	case MeditHealthInput:
		handleMeditHealthInput(player)
	case MeditMana:
		handleMeditManaMenu(player)
	case MeditManaInput:
		handleMeditManaInput(player)
	case MeditMove:
		handleMeditMoveMenu(player)
	case MeditMoveInput:
		handleMeditMoveInput(player)
	case MeditSave:
		handleMeditSave(player)
	}
}

func handleMeditSave(player *Player) {

	if player.Olc.Medit.EditMob.Id == -1 {
		_, err := player.Olc.Medit.EditMob.Zone.AddNewMob(player.Olc.Medit.EditMob)
		if err != nil {
			log.Printf("Error adding mob to zone: %v\n", err)
			player.SendToCharacter("Error adding mob to zone")
		}
	} else {
		player.Olc.Medit.EditMob.Zone.Mobs[player.Olc.Medit.EditMob.Id] = player.Olc.Medit.EditMob
	}

	err := player.Olc.Medit.EditMob.Save()
	if err != nil {
		log.Printf("Error saving mob: %v\n", err)
		player.SendToCharacter("Failed to save Mob")

	} else {
		player.SendToCharacter("Mob saved.\r\n")
	}
	player.State = PlayerStateActive
	player.Olc.Mode = Inactive
	player.Olc.Medit = nil

}
func handleMeditMainMenu(player *Player) {
	player.SendToCharacter(clearScreenCommand)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("(1)Name: %s\r\n", player.Olc.Medit.EditMob.Name))
	builder.WriteString(fmt.Sprintf("(2)Description:\r\n%s\r\n", player.Olc.Medit.EditMob.Description))
	builder.WriteString(fmt.Sprintf("(3)Health: %d\r\n", player.Olc.Medit.EditMob.Resources.MaxHealth))
	builder.WriteString(fmt.Sprintf("(4)Mana: %d\r\n", player.Olc.Medit.EditMob.Resources.MaxMana))
	builder.WriteString(fmt.Sprintf("(5)Move: %d\r\n", player.Olc.Medit.EditMob.Resources.MaxMove))
	builder.WriteString("(A)Abort\r\n")
	builder.WriteString("(S)Save\r\n")
	builder.WriteString("Selection: ")
	player.SendToCharacter(builder.String())
	player.Olc.Medit.State = MeditMainMenuSelection
}

func handleMeditMainMenuSelection(player *Player) {
	option, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	option = strings.ToLower(strings.TrimSpace(option))
	switch option {
	case "1":
		player.Olc.Medit.State = MeditNameMenu
	case "2":
		player.Olc.Medit.State = MeditDescription
	case "3":
		player.Olc.Medit.State = MeditHealth
	case "4":
		player.Olc.Medit.State = MeditMana
	case "5":
		player.Olc.Medit.State = MeditMove
	case "a":
		player.State = PlayerStateActive
		player.Olc.Mode = Inactive
		player.Olc.Medit = nil
	case "s":
		player.Olc.Medit.State = MeditSave
	}
}

func handleMeditNameMenu(player *Player) {
	player.SendToCharacter(clearScreenCommand)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Name: %s\r\n", player.Olc.Medit.EditMob.Name))
	builder.WriteString("Enter new name: ")
	player.SendToCharacter(builder.String())
	player.Olc.Medit.State = MeditNameInput
}

func handleMeditNameInput(player *Player) {
	name, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	name = strings.TrimSpace(name)
	player.Olc.Medit.EditMob.Name = name
	player.Olc.Medit.State = MeditMainMenu
}

func handleMeditDescription(player *Player) {
	builder := strings.Builder{}
	builder.WriteString("Instructions: /s to save, /h for more options.\r\n")
	for _, line := range player.Olc.Medit.DescriptionBuffer {
		builder.WriteString(line)
		builder.WriteString("\r\n")
	}
	builder.WriteString("> ")
	player.SendToCharacter(builder.String())
	player.Olc.Medit.State = MeditDescriptionInput
}

func handleMeditDescriptionInput(player *Player) {
	description, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	description = strings.TrimSpace(description)

	if len(description) == 2 {
		switch description {
		case "/s":
			player.Olc.Medit.EditMob.Description = strings.Join(player.Olc.Medit.DescriptionBuffer, "\r\n")
			player.Olc.Medit.State = MeditMainMenu
			return
		case "/h":
			player.SendToCharacter("Instructions: /s to save, /h for more options.\r\n")
			return
		}
	}

	player.Olc.Medit.DescriptionBuffer = append(player.Olc.Medit.DescriptionBuffer, description)
}

func handleMeditHealthMenu(player *Player) {
	player.SendToCharacter(clearScreenCommand)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Health: %d\r\n", player.Olc.Medit.EditMob.Resources.MaxHealth))
	builder.WriteString("Enter new health: ")
	player.SendToCharacter(builder.String())
	player.Olc.Medit.State = MeditHealthInput
}

func handleMeditHealthInput(player *Player) {
	input, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	input = strings.TrimSpace(input)
	health, err := strconv.Atoi(input)
	if err != nil || health <= 0 {
		player.SendToCharacter("Must be a number greater then 0")
		player.Olc.Medit.State = MeditHealth
	}
	player.Olc.Medit.EditMob.Resources.MaxHealth = health
	player.Olc.Medit.State = MeditMainMenu
}

func handleMeditManaMenu(player *Player) {
	player.SendToCharacter(clearScreenCommand)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Mana: %d\r\n", player.Olc.Medit.EditMob.Resources.MaxMana))
	builder.WriteString("Enter new Mana: ")
	player.SendToCharacter(builder.String())
	player.Olc.Medit.State = MeditManaInput
}

func handleMeditManaInput(player *Player) {
	input, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	input = strings.TrimSpace(input)
	mana, err := strconv.Atoi(input)
	if err != nil || mana <= 0 {
		player.SendToCharacter("Must be a number greater then 0")
		player.Olc.Medit.State = MeditMana
	}
	player.Olc.Medit.EditMob.Resources.MaxMana = mana
	player.Olc.Medit.State = MeditMainMenu
}

func handleMeditMoveMenu(player *Player) {
	player.SendToCharacter(clearScreenCommand)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("Move: %d\r\n", player.Olc.Medit.EditMob.Resources.MaxMove))
	builder.WriteString("Enter new move: ")
	player.SendToCharacter(builder.String())
	player.Olc.Medit.State = MeditMoveInput
}

func handleMeditMoveInput(player *Player) {
	input, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	input = strings.TrimSpace(input)
	move, err := strconv.Atoi(input)
	if err != nil || move <= 0 {
		player.SendToCharacter("Must be a number greater then 0")
		player.Olc.Medit.State = MeditMove
	}
	player.Olc.Medit.EditMob.Resources.MaxMove = move
	player.Olc.Medit.State = MeditMainMenu
}
