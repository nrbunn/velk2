package velk

type Zone struct {
	Id    int
	Name  string
	Rooms []*Room
}

func NewZone(id int) *Zone {
	return &Zone{
		Id:    id,
		Name:  "Untitled Zone",
		Rooms: make([]*Room, 0),
	}
}

func (z *Zone) SetName(name string) {
	z.Name = name
}
func (z *Zone) CreateRoom() *Room {
	room := NewRoom(len(z.Rooms), z)
	z.Rooms = append(z.Rooms, room)
	return room
}

func (z *Zone) GetRoom(id int) *Room {
	return z.Rooms[id]
}
