package structs

import (
	"github.com/dropbox/godropbox/container/bitvector"
	"github.com/google/uuid"
	"log"
	"velk2/pkg/interfaces"
	"velk2/pkg/libs"
)

const PlayerStateLoading = "PLAYER_STATE_LOADING"
const PlayerStateActive = "PLAYER_STATE_ACTIVE"
const PlayerStateNaming = "PLAYER_STATE_NAMING"

type Player struct {
	UUID          uuid.UUID
	Connection    interfaces.ConnectionInterface
	Name          string
	CommandBuffer *libs.Queue
	Server        interfaces.ServerInterface
	Position      bitvector.BitVector
	State         string
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

func (p *Player) GetUUID() string {
	return uuid.NewString()
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

		}
		p.CommandBuffer.Insert(string(data))
		log.Println("Server got:", string(data))
	}
}
