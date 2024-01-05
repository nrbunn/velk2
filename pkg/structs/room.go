package structs

import (
	"github.com/wk8/go-ordered-map/v2"
	"velk2/pkg/interfaces"
)

var Directions = []string{"north", "east", "south", "west", "up", "down"}

type Room struct {
	Id          int
	Name        string
	Description string
	Characters  []interfaces.Character
	Exits       *orderedmap.OrderedMap[string, interfaces.Exit]
	//Objects []*interfaces.Objects

}

func NewRoom(id int) *Room {

	exits := orderedmap.New[string, interfaces.Exit]()
	for _, direction := range Directions {
		exits.Set(direction, nil)
	}
	return &Room{
		Id:          id,
		Name:        "Untitled Room",
		Description: "This room has no description",
		Characters:  make([]interfaces.Character, 0),
		Exits:       exits,
	}
}

func (r *Room) GetName() string {
	return r.Name
}

func (r *Room) GetDescription() string {
	return r.Description
}

func (r *Room) GetExits() string {
	message := "[ "
	for pair := r.Exits.Oldest(); pair != nil; pair = pair.Next() {
		if pair.Value != nil {
			message += pair.Key + " "
		}
	}
	message += "]"
	return message
}

func (r *Room) GetCharacters() []interfaces.Character {
	return r.Characters
}

func (r *Room) AddCharacter(targetCharacter interfaces.Character) {
	r.Characters = append(r.Characters, targetCharacter)
}

func (r *Room) RemoveCharacter(targetCharacter interfaces.Character) {
	for i, character := range r.Characters {
		if character == targetCharacter {
			r.Characters = append(r.Characters[:i], r.Characters[i+1:]...)
		}
	}
}

func (r *Room) SetExit(dir string, room interfaces.Room) error {

	exit := &Exit{
		Room: room,
	}

	r.Exits.Set(dir, exit)

	return nil
}

func (r *Room) GetExit(dir string) interfaces.Exit {

	if exit, ok := r.Exits.Get(dir); ok {
		return exit
	}

	return nil
}

func (r *Room) SetName(name string) {
	r.Name = name
}
