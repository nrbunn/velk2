package interfaces

type Character interface {
	GetUUID() string
	GetName() string
	SendToCharacter(string) error
	GetRoom() Room
	SetRoom(Room)
}
