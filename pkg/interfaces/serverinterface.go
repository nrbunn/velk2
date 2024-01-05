package interfaces

type ServerInterface interface {
	GetPlayers() []Character
	GetColorService() ColorServiceInterface
}
