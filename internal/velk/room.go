package velk

import (
	"encoding/json"
	"fmt"
	"github.com/wk8/go-ordered-map/v2"
	"os"
	"path/filepath"
	"strconv"
	"velk2/internal/utils"
)

type Room struct {
	Id          int                                  `json:"id"`
	Name        string                               `json:"name"`
	Description string                               `json:"description"`
	Characters  []*Player                            `json:"-"`
	Mobs        []*Mob                               `json:"-"`
	Exits       *orderedmap.OrderedMap[string, Rnum] `json:"exits"`
	Zone        *Zone                                `json:"-"`
	//Objects []*interfaces.Objects

}

func NewRoom(id int, zone *Zone) *Room {
	exits := orderedmap.New[string, Rnum]()
	for _, direction := range utils.Directions {
		exits.Set(direction, Rnum{ZoneId: -1, RoomId: -1})
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

func (r *Room) Save() error {
	dir := filepath.Join("wld", "zone", strconv.Itoa(r.Zone.Id))
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(dir, fmt.Sprintf("%d.json", r.Id)))
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (r *Room) GetRnum() Rnum {
	return Rnum{ZoneId: r.Zone.Id, RoomId: r.Id}
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

func (r *Room) SetExit(dir string, rnum Rnum) error {

	r.Exits.Set(dir, rnum)

	return nil
}

func (r *Room) GetExit(dir string) Rnum {

	if exit, ok := r.Exits.Get(dir); ok {
		return exit
	}

	return Rnum{ZoneId: -1, RoomId: -1}
}
