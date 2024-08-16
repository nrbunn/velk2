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

type SpawnType int
type SpawnStrategy int

const (
	SpawnTypeMob    SpawnType = 1
	SpawnTypeObject SpawnType = 2
)
const (
	SpawnStrategyRoomSpawnCount SpawnStrategy = 1
)

type SpawnData struct {
	Vnum      Vnum
	SpawnType SpawnType
}

type Room struct {
	Id               int                                  `json:"id"`
	Name             string                               `json:"name"`
	Description      string                               `json:"description"`
	Characters       []*Player                            `json:"-"`
	Mobs             []*Mob                               `json:"-"`
	Exits            *orderedmap.OrderedMap[string, Vnum] `json:"exits"`
	Zone             *Zone                                `json:"-"`
	SpawnInformation []SpawnData                          `json:"spawn_data"`

	//Objects []*interfaces.Objects

}

func NewRoom(id int, zone *Zone) *Room {
	exits := orderedmap.New[string, Vnum]()
	for _, direction := range utils.Directions {
		exits.Set(direction, Vnum{ZoneId: -1, VirtualId: -1})
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

func (r *Room) SpawnMobs(z *Zone) {
	for _, spawnData := range r.SpawnInformation {
		if len(r.Mobs) >= r.MaxMobs || mobExistsInZone(z, spawnData.Mob) {
			continue
		}

		newMob := NewMob(spawnData.Mob.Name, spawnData.Mob.HP)
		r.Mobs[newMob.ID] = newMob
		z.Mobs[newMob.ID] = newMob
		fmt.Printf("A %s has spawned in room %s\n", newMob.Name, r.Name)
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

func (r *Room) GetRnum() Vnum {
	return Vnum{ZoneId: r.Zone.Id, VirtualId: r.Id}
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

func (r *Room) SetExit(dir string, rnum Vnum) error {

	r.Exits.Set(dir, rnum)

	return nil
}

func (r *Room) GetExit(dir string) Vnum {

	if exit, ok := r.Exits.Get(dir); ok {
		return exit
	}

	return Vnum{ZoneId: -1, VirtualId: -1}
}
