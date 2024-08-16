package velk

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"velk2/internal/utils"
)

type VelkServer struct {
	Server net.Listener
	Addr   string
}

func NewVelkServer(addr string) (*VelkServer, error) {
	server, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &VelkServer{
		Server: server,
		Addr:   addr,
	}, nil
}

func initializeCommands() *Commands {
	cmds := NewCommands()

	for _, direction := range utils.Directions {
		cmds.AddCommand(direction, Move)
	}
	cmds.AddCommand("gossip", Gossip)
	cmds.AddCommand("kill", Kill)
	cmds.AddCommand("look", Look)
	cmds.AddCommand("medit", Medit)
	cmds.AddCommand("redit", Redit)
	cmds.AddCommand("say", Say)
	return cmds
}

func initializeRooms() {
	err := LoadZones()
	if err != nil {
		log.Fatal(err)
	}

	//mob := NewMob()
	//mob.Name = "goblin"
	//mob.ShortDesc = "a goblin"
	//mob.Resources.Health = 5
	//mob.Room = room
	//room.AddMob(mob)

}

func (v *VelkServer) GameLoop() {
	tickChan := make(chan time.Time)
	newPlayerChan := make(chan *Player)
	Cmds := initializeCommands()
	initializeRooms()

	go v.acceptConnections(newPlayerChan)
	go v.tick(tickChan)
	for {
		select {
		case player := <-newPlayerChan:
			Players = append(Players, player)
		case <-tickChan:
			v.executePlayerActions(Cmds)
			pulseViloence()
		}
	}
}

func pulseViloence() {
	for _, player := range Players {
		if len(player.FightingTargets) != 0 {
			targetMob := player.FightingTargets[0]
			if !isReadyToAttack(player) {
				continue
			}

			player.SendToCharacter(fmt.Sprintf("You hit %s!\r\n", targetMob.Name))
			player.LastAttack = time.Now()
			targetMob.Resources.Health -= 1
			if targetMob.Resources.Health <= 0 {
				player.SendToCharacter(fmt.Sprintf("%s is dead!\r\n", targetMob.Name))
				targetMob.Room.RemoveMob(targetMob)
				player.RemoveFightingTarget(targetMob)
			}
		}
	}
}

func isReadyToAttack(player *Player) bool {
	attackSpeedDuration := time.Duration(player.AttackSpeed * float32(time.Second))
	return time.Since(player.LastAttack) >= attackSpeedDuration
}

func (v *VelkServer) executePlayerActions(cmds *Commands) {
	for _, player := range Players {
		switch player.State {
		case PlayerStateLoading:
			handlePlayerStateLoading(player)
		case PlayerStateName:
			handlePlayerStateName(player)
		case PlayerStateActive:
			handlePlayerStateActive(player, cmds)
		case PlayerStateEditing:
			HandleOlc(player)
		}
	}
}

// handlePlayerStateLoading handles the player's loading state
func handlePlayerStateLoading(player *Player) {
	player.SendToCharacter("What is thy name?")
	player.State = PlayerStateName
}

// handlePlayerStateName handles the player's name state
func handlePlayerStateName(player *Player) {

	out, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	out = sanitizeOutput(out)
	if utils.HasSpecialChar(out) {
		player.SendToCharacter("Name cannot have special characters.\r\n")
		player.State = PlayerStateLoading
		return
	}
	player.Load(out)
	player.Save()
	player.SendToCharacter(fmt.Sprintf("Welcome %s\r\n", player.Name))
	player.State = PlayerStateActive
	room := Zones[0].Rooms[0]
	player.Room = room
	room.AddCharacter(player)
}

// handlePlayerStateActive handles the player's active state
func handlePlayerStateActive(player *Player, cmds *Commands) {
	out, err := player.CommandBuffer.Remove()
	if err != nil {
		return
	}
	cmd, cmdSuffix := processCommandString(out)

	action := cmds.GetCommand(cmd)
	if action != nil {
		action(player, cmd, cmdSuffix)
	} else {
		player.SendToCharacter("Huh?\r\n")
	}
}

// sanitizeOutput sanitizes the command output string
func sanitizeOutput(out string) string {
	return strings.ToLower(strings.TrimSpace(out))
}

func processCommandString(commandString string) (string, string) {

	commandSplit := strings.SplitN(commandString, " ", 2)

	if len(commandSplit) < 2 {
		return strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(commandSplit[0], "\n"), "\r")), ""
	}

	return strings.ToLower(commandSplit[0]), strings.TrimSuffix(strings.TrimSuffix(commandSplit[1], "\n"), "\r")
}

func (v *VelkServer) acceptConnections(newConn chan *Player) {
	for {
		conn, err := v.Server.Accept()
		if err != nil {
			log.Println("Failed to accept conn.", err)
			continue
		}

		player, err := NewPlayer(NewConnection(conn))
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
