package interfaces

type ConnectionInterface interface {
	Read() (string, error)
	Write(s string) error
	Close()
}
