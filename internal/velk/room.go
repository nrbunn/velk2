package velk

import (
	"github.com/wk8/go-ordered-map/v2"
	"velk2/internal/utils"
)

type Room struct {
	Id          int
	Name        string
	Description string
	Characters  []*Player
	Mobs        []*Mob
	Exits       *orderedmap.OrderedMap[string, *Exit]
	Zone        *Zone
	//Objects []*interfaces.Objects

}

func NewRoom(id int, zone *Zone) *Room {
	exits := orderedmap.New[string, *Exit]()
	for _, direction := range utils.Directions {
		exits.Set(direction, nil)
	}
	return &Room{
		Id:          id,
		Name:        "Untitled Room",
		Description: "This room has no description",
		Characters:  make([]*Player, 0),
		Mobs:        make([]*Mob, 0),
		Exits:       exits,
		Zone:        zone,
	}
}

func (r *Room) AddCharacter(targetCharacter *Player) {
	r.Characters = append(r.Characters, targetCharacter)
}

func (r *Room) RemoveCharacter(targetCharacter *Player) {
	for i, character := range r.Characters {
		if character == targetCharacter {
			r.Characters = append(r.Characters[:i], r.Characters[i+1:]...)
		}
	}
}

func (r *Room) AddMob(targetMob *Mob) {
	r.Mobs = append(r.Mobs, targetMob)
}

func (r *Room) RemoveMob(targetMob *Mob) {
	for i, mob := range r.Mobs {
		if mob == targetMob {
			r.Mobs = append(r.Mobs[:i], r.Mobs[i+1:]...)
		}
	}
}

func (r *Room) SetExit(dir string, room *Room) error {

	exit := &Exit{
		Room: room,
	}

	r.Exits.Set(dir, exit)

	return nil
}

func (r *Room) GetExit(dir string) *Exit {

	if exit, ok := r.Exits.Get(dir); ok {
		return exit
	}

	return nil
}
