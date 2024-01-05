package structs

import "velk2/pkg/interfaces"

type Zone struct {
	Id    int
	Name  string
	Rooms []interfaces.Room
}

func NewZone(id int) Zone {
	return Zone{
		Id:    id,
		Name:  "Untitled Zone",
		Rooms: make([]interfaces.Room, 0),
	}
}

func (z *Zone) SetName(name string) {
	z.Name = name
}
func (z *Zone) CreateRoom() interfaces.Room {
	room := NewRoom(len(z.Rooms) + 1)
	z.Rooms = append(z.Rooms, room)
	return room
}

func (z *Zone) GetRoom(id int) interfaces.Room {
	return z.Rooms[id]
}
