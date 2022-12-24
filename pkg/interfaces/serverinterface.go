package interfaces

type ServerInterface interface {
	GetPlayers() []CharacterInterface
	GetColorService() ColorServiceInterface
}
