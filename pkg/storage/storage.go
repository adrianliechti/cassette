package storage

type Event any

type Storage interface {
	Events(session string) ([]Event, error)

	AppendEvents(session string, events ...Event) error
	DeleteDelete(session string) error
}
