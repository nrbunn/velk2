package interfaces

type CharacterInterface interface {
	GetUUID() string
	SendToCharacter(string) error
}
