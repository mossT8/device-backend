package store

type Store interface {
	Close() (error, error)
	Ping() error
}
