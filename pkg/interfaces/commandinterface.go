package interfaces

type CommandInterface interface {
	Action(CharacterInterface, string, ...string)
}
