package repository

import (
	"time"
)

type Session struct {
	ID string `json:"id"`

	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`

	Origin string

	User      string
	UserName  string
	UserAgent string
}

type Repository interface {
	Sessions() ([]Session, error)
	Session(id string) (*Session, error)

	CreateSession(info *SessionInfo) (*Session, error)
	DeleteSession(id string) error
}

type SessionInfo struct {
	Origin string

	User      string
	UserName  string
	UserAgent string
}
