package velk

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strconv"
)

type Mob struct {
	Id              int        `json:"id"`
	Uuid            uuid.UUID  `json:"-"`
	Name            string     `json:"name"`
	ShortDesc       string     `json:"shortdesc"`
	Description     string     `json:"description"`
	Room            *Room      `json:"-"`
	Resources       *Resources `json:"resources"`
	Zone            *Zone      `json:"-"`
	FightingTargets []*Player  `json:"-"`
}

func NewMob(id int, z *Zone) *Mob {
	resources := &Resources{
		Health:    100,
		MaxHealth: 100,
		Mana:      100,
		MaxMana:   100,
		Move:      100,
		MaxMove:   100,
	}

	return &Mob{
		Id:              id,
		Name:            "Untitled Mob",
		ShortDesc:       "An untitled mob",
		Resources:       resources,
		Zone:            z,
		FightingTargets: make([]*Player, 0),
	}
}

func (r *Mob) Save() error {
	dir := filepath.Join("wld", "zone", strconv.Itoa(r.Zone.Id), "mobs")
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

func (r *Mob) SpawnMob() *Mob {
	m := *r
	m.Uuid = uuid.New()
	return &m
}
