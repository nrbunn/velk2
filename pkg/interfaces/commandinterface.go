package interfaces

type CommandInterface interface {
	Action(Character, string, ...string)
}
