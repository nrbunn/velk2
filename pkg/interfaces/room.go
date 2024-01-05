package interfaces

type Room interface {
	GetName() string
	GetCharacters() []Character
	RemoveCharacter(targetCharacter Character)
	AddCharacter(targetCharacter Character)
	GetDescription() string
	GetExits() string
	GetExit(string) Exit
	SetExit(string, Room) error
	SetName(string)
}
