package datastore

type DataStore interface {
	Start() error
	Close() error
	Ping() error
}
