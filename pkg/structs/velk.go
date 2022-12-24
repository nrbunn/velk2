package structs

import (
	"log"
	"net"
	"strings"
	"time"
	"velk2/pkg/commands"
	"velk2/pkg/interfaces"
)

type VelkServer struct {
	Server      net.Listener
	Addr        string
	Players     []interfaces.CharacterInterface
	ColorServce interfaces.ColorServiceInterface
}

func NewVelkServer(addr string) (*VelkServer, error) {
	server, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &VelkServer{
		Server:      server,
		Addr:        addr,
		Players:     make([]interfaces.CharacterInterface, 0),
		ColorServce: NewColorService(),
	}, nil
}

func (v *VelkServer) GetPlayers() []interfaces.CharacterInterface {
	return v.Players
}

func (v *VelkServer) GetColorService() interfaces.ColorServiceInterface {
	return v.ColorServce
}

func (v *VelkServer) GameLoop() {
	tickChan := make(chan time.Time)
	newConnChan := make(chan *Player)
	CommandList := make(map[string]interfaces.CommandInterface)
	CommandList["say"] = &commands.Say{Server: v}

	go v.acceptConnections(newConnChan)
	go v.tick(tickChan)
	for {
		select {
		case player := <-newConnChan:
			v.Players = append(v.Players, player)
		case <-tickChan:
			for _, playerInterface := range v.Players {
				player := playerInterface.(*Player)
				out, err := player.CommandBuffer.Remove()
				if err != nil {
					continue
				}
				cmd, cmdSuffix := processCommandString(out)
				action, actionExists := CommandList[cmd]
				if actionExists {
					action.Action(player, cmd, cmdSuffix)
				} else {
					player.SendToCharacter("Huh?\r\n")
				}
			}
		}
	}
}

func processCommandString(commandString string) (string, string) {

	commandSplit := strings.SplitN(commandString, " ", 2)

	if len(commandSplit) < 2 {
		return commandSplit[0], ""
	}

	return strings.ToLower(commandSplit[0]), commandSplit[1]
}

func (v *VelkServer) acceptConnections(newConn chan *Player) {
	for {
		conn, err := v.Server.Accept()
		if err != nil {
			log.Println("Failed to accept conn.", err)
			continue
		}

		player, err := NewPlayer(conn, v)
		if err != nil {
			log.Println("failed to create player.", err)
			continue
		}
		newConn <- player
	}
}

func (v *VelkServer) tick(tick chan time.Time) {
	for {
		time.Sleep(100 * time.Microsecond)
		tick <- time.Now()
	}
}
