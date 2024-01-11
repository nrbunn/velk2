package velk

type Exit struct {
	Room *Room
}

func NewExit(room *Room) *Exit {
	return &Exit{Room: room}
}

func (e *Exit) TargetRoom() *Room {
	return e.Room
}
