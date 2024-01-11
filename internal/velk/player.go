package velk

import (
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/container/bitvector"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
	"velk2/internal/libs"
)

const PlayerStateLoading = "PLAYER_STATE_LOADING"
const PlayerStateActive = "PLAYER_STATE_ACTIVE"
const PlayerStateName = "PLAYER_STATE_NAME"
const PlayerStateEditing = "PLAYER_STATE_OLC"

type Resources struct {
	Health    int `json:"health"`
	MaxHealth int `json:"maxhealth"`
	Mana      int `json:"mana"`
	MaxMana   int `json:"maxmana"`
	Move      int `json:"move"`
	MaxMove   int `json:"maxmove"`
}

type Player struct {
	UUID            uuid.UUID           `json:"uuid"`
	Connection      *Connection         `json:"-"`
	Name            string              `json:"name"`
	CommandBuffer   *libs.Queue         `json:"-"`
	Position        bitvector.BitVector `json:"position"`
	State           string              `json:"-"`
	Room            *Room               `json:"-"`
	RoomId          string              `json:"roomid"`
	Resources       *Resources          `json:"resources"`
	FightingTargets []*Mob              `json:"-"`
	AttackSpeed     float32             `json:"attackspeed"`
	LastAttack      time.Time           `json:"-"`
	Olc             *Olc                `json:"-"`
}

func NewPlayer(connection *Connection) (*Player, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	resources := &Resources{
		Health:    100,
		MaxHealth: 100,
		Mana:      100,
		MaxMana:   100,
		Move:      100,
		MaxMove:   100,
	}

	player := &Player{
		UUID:            uuid,
		Connection:      connection,
		Name:            "",
		CommandBuffer:   libs.NewQueue(300),
		State:           PlayerStateLoading,
		Resources:       resources,
		FightingTargets: make([]*Mob, 0),
		AttackSpeed:     3.0,
		Olc:             &Olc{Mode: Inactive, Redit: nil},
	}

	go player.readConn()
	return player, nil
}

func (p *Player) Load(name string) {
	log.Println("Loading")
	p.Name = name

	fn := filepath.Join("wld", "plr", string(p.Name[0]), fmt.Sprintf("%s.json", p.Name))
	if _, err := os.Stat(fn); err != nil {
		log.Println("does not exists")
		return
	}

	f, err := os.Open(fn)
	if err != nil {
		log.Println(fmt.Sprintf("failed to open player file %s for %s", fn, p.Name), err)
		return
	}
	defer f.Close()

	b, _ := ioutil.ReadAll(f)

	err = json.Unmarshal(b, p)
	if err != nil {
		log.Println(fmt.Sprintf("failed to unmarshal player %s", p.Name), err)
	}
}

func (p *Player) Save() {
	log.Println("Saving")

	b, err := json.Marshal(p)
	if err != nil {
		log.Println(fmt.Sprintf("failed to marshal %s", p.Name), err)
		return
	}

	fn := filepath.Join("wld", "plr", string(p.Name[0]), fmt.Sprintf("%s.json", p.Name))
	f, err := os.Create(fn)
	if err != nil {
		log.Println(fmt.Sprintf("failed to create file %s for %s", fn, p.Name), err)
		return
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		log.Println(fmt.Sprintf("failed to write file %s", p.Name), err)
	}
}

func (p *Player) AddFightingTarget(targetMob *Mob) {
	p.FightingTargets = append(p.FightingTargets, targetMob)
}

func (p *Player) RemoveFightingTarget(targetMob *Mob) {
	for i, mob := range p.FightingTargets {
		if mob == targetMob {
			p.FightingTargets = append(p.FightingTargets[:i], p.FightingTargets[i+1:]...)
		}
	}
}

func (p *Player) SendToCharacter(data string) error {
	data = libs.ProcessString(data)
	err := p.Connection.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (p *Player) closeConnection() {
	p.Connection.Close()
}

func (p *Player) readConn() {
	for {
		data, err := p.Connection.Read()
		if err != nil {
			break
		}
		p.CommandBuffer.Insert(data)
		log.Println("Server got:", data)
	}
}
