package structs

import (
	"github.com/google/uuid"
	"log"
	"net"
	"velk2/pkg/interfaces"
	"velk2/pkg/libs"
)

type Player struct {
	UUID          uuid.UUID
	Connection    net.Conn
	Name          string
	CommandBuffer *libs.Queue
	Server        interfaces.ServerInterface
}

func NewPlayer(conn net.Conn, server interfaces.ServerInterface) (*Player, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	player := &Player{
		UUID:          uuid,
		Connection:    conn,
		Name:          "",
		CommandBuffer: libs.NewQueue(300),
		Server:        server,
	}
	go player.readConn()
	return player, nil
}

func (p *Player) GetUUID() string {
	return uuid.NewString()
}

func (p *Player) SendToCharacter(data string) error {
	data = p.Server.GetColorService().ProcessString(data)
	_, err := p.Connection.Write([]byte(data))
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
		buf := make([]byte, 512)
		nr, err := p.Connection.Read(buf)
		if err != nil {
			p.closeConnection()
			return
		}

		data := buf[0 : nr-1]
		p.CommandBuffer.Insert(string(data))
		log.Println("Server got:", string(data))
	}
}
