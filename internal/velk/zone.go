package velk

import (
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
)

type Zone struct {
	Id         int
	Name       string
	Rooms      map[int]*Room      `json:"-"`
	Mobs       map[int]*Mob       `json:"-"`
	ActiveMobs map[uuid.UUID]*Mob `json:"-"`
}

func NewZone(id int) *Zone {
	return &Zone{
		Id:    id,
		Name:  "Untitled Zone",
		Rooms: make(map[int]*Room),
	}
}

func (z *Zone) SetName(name string) {
	z.Name = name
}

func (z *Zone) CreateRoom() *Room {
	room := NewRoom(len(z.Rooms), z)
	z.Rooms[room.Id] = room
	return room
}

func (z *Zone) AddNewMob(mob *Mob) (*Mob, error) {
	mob.Id = len(z.Mobs)
	z.Mobs[mob.Id] = mob
	return mob, nil
}

func (z *Zone) GetRoom(id int) *Room {
	return z.Rooms[id]
}

func (z *Zone) SaveZone() error {
	data, err := json.Marshal(z)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join("wld", "zone", strconv.Itoa(z.Id), "metadata.json"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (z *Zone) LoadRooms() error {
	dir := filepath.Join("wld", "zone", strconv.Itoa(z.Id), "rooms")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			roomFile := filepath.Join(dir, file.Name())
			log.Println("Loading room", roomFile)
			data, err := ioutil.ReadFile(roomFile)
			if err != nil {
				return err
			}
			room := &Room{}
			err = json.Unmarshal(data, room)
			if err != nil {
				return err
			}
			room.Zone = z
			z.Rooms[room.Id] = room
		}
	}
	return nil
}

func (z *Zone) LoadMobs() error {
	dir := filepath.Join("wld", "zone", strconv.Itoa(z.Id), "mobs")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			mobFile := filepath.Join(dir, file.Name())
			log.Println("Loading mob", mobFile)
			data, err := ioutil.ReadFile(mobFile)
			if err != nil {
				return err
			}
			mob := &Mob{}
			err = json.Unmarshal(data, mob)
			if err != nil {
				return err
			}
			mob.Zone = z
			z.Mobs[mob.Id] = mob
		}
	}
	return nil
}
