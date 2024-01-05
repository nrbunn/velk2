package structs

import (
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/container/bitvector"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"velk2/pkg/interfaces"
	"velk2/pkg/libs"
)

const PlayerStateLoading = "PLAYER_STATE_LOADING"
const PlayerStateActive = "PLAYER_STATE_ACTIVE"
const PlayerStateName = "PLAYER_STATE_NAME"

type Player struct {
	UUID          uuid.UUID                      `json:"uuid"`
	Connection    interfaces.ConnectionInterface `json:"-"`
	Name          string                         `json:"name"`
	CommandBuffer *libs.Queue                    `json:"-"`
	Server        interfaces.ServerInterface     `json:"-"`
	Position      bitvector.BitVector            `json:"position"`
	State         string                         `json:"-"`
	Room          interfaces.Room                `json:"-"`
	RoomId        string                         `json:"roomid"`
}

func NewPlayer(connection interfaces.ConnectionInterface, server interfaces.ServerInterface) (*Player, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	player := &Player{
		UUID:          uuid,
		Connection:    connection,
		Name:          "",
		CommandBuffer: libs.NewQueue(300),
		Server:        server,
		State:         PlayerStateLoading,
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
func (p *Player) GetUUID() string {
	return p.UUID.String()
}

func (p *Player) GetName() string {

	return strings.ToTitle(p.Name)
}

func (p *Player) GetRoom() interfaces.Room {
	return p.Room
}

func (p *Player) SetRoom(room interfaces.Room) {
	p.Room = room
}

func (p *Player) SendToCharacter(data string) error {
	data = p.Server.GetColorService().ProcessString(data)
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
