package repository

import (
	"time"
)

type Session struct {
	ID string `json:"id"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type Repository interface {
	Sessions() ([]Session, error)
	Session(id string) (*Session, error)

	CreateSession() (*Session, error)
	DeleteSession(id string) error
}
