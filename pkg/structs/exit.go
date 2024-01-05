package structs

import "velk2/pkg/interfaces"

type Exit struct {
	Room interfaces.Room
}

func NewExit(room interfaces.Room) *Exit {
	return &Exit{Room: room}
}

func (e *Exit) TargetRoom() interfaces.Room {
	return e.Room
}
