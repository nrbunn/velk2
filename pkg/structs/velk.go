package structs

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"velk2/pkg/commands"
	"velk2/pkg/interfaces"
	"velk2/pkg/utils"
)

type VelkServer struct {
	Server       net.Listener
	Addr         string
	Players      []interfaces.Character
	Connections  []interfaces.ConnectionInterface
	ColorService interfaces.ColorServiceInterface
}

func NewVelkServer(addr string) (*VelkServer, error) {
	server, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &VelkServer{
		Server:       server,
		Addr:         addr,
		Players:      make([]interfaces.Character, 0),
		ColorService: NewColorService(),
	}, nil
}

func (v *VelkServer) GetPlayers() []interfaces.Character {
	return v.Players
}

func (v *VelkServer) GetColorService() interfaces.ColorServiceInterface {
	return v.ColorService
}

func (v *VelkServer) GameLoop() {
	tickChan := make(chan time.Time)
	newPlayerChan := make(chan *Player)
	Cmds := NewCommands()
	Cmds.addCommand("north", &commands.Move{Server: v})
	Cmds.addCommand("east", &commands.Move{Server: v})
	Cmds.addCommand("south", &commands.Move{Server: v})
	Cmds.addCommand("west", &commands.Move{Server: v})
	Cmds.addCommand("up", &commands.Move{Server: v})
	Cmds.addCommand("down", &commands.Move{Server: v})
	Cmds.addCommand("gossip", &commands.Gossip{Server: v})
	Cmds.addCommand("look", &commands.Look{Server: v})
	Cmds.addCommand("say", &commands.Say{Server: v})

	zone := NewZone(1)
	room := zone.CreateRoom()
	room.SetName("The Void")
	room2 := zone.CreateRoom()

	room.SetExit("north", room2)
	room2.SetExit("south", room)

	go v.acceptConnections(newPlayerChan)
	go v.tick(tickChan)
	for {
		select {
		case player := <-newPlayerChan:
			v.Players = append(v.Players, player)
		case <-tickChan:
			for _, playerInterface := range v.Players {
				player := playerInterface.(*Player)
				switch player.State {
				case PlayerStateLoading:
					{
						player.SendToCharacter("What is thy name?")
						player.State = PlayerStateName
					}
				case PlayerStateName:
					{
						out, err := player.CommandBuffer.Remove()
						if err != nil {
							continue
						}
						out = strings.ToLower(strings.TrimSpace(out))

						if utils.HasSpecialChar(out) {
							player.SendToCharacter("Name cannot have special characters.\r\n")
							player.State = PlayerStateLoading
							break
						}

						//todo check name length

						player.Load(strings.TrimSpace(out))
						player.Save()
						player.SendToCharacter(fmt.Sprintf("Welcome %s\r\n", player.GetName()))
						player.State = PlayerStateActive
						player.Room = room
						room.AddCharacter(player)
					}

				case PlayerStateActive:
					{
						out, err := player.CommandBuffer.Remove()
						if err != nil {
							continue
						}
						cmd, cmdSuffix := processCommandString(out)
						log.Println(fmt.Sprintf("command:'%s'", cmd))
						action := Cmds.getCommand(cmd)
						if action != nil {
							action.Action(player, cmd, cmdSuffix)
						} else {
							player.SendToCharacter("Huh?\r\n")
						}
					}
				}
			}
		}
	}
}

func processCommandString(commandString string) (string, string) {

	commandSplit := strings.SplitN(commandString, " ", 2)

	if len(commandSplit) < 2 {
		return strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(commandSplit[0], "\n"), "\r")), ""
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

		player, err := NewPlayer(NewConnection(conn), v)
		if err != nil {
			log.Println("failed to create player.", err)
			continue
		}
		newConn <- player
	}
}

func (v *VelkServer) tick(tick chan time.Time) {
	ticker := time.NewTicker(100 * time.Microsecond)
	defer ticker.Stop()
	for {
		tick <- <-ticker.C
	}
}
