package config

import (
	"cassette/pkg/repository"
	"cassette/pkg/repository/gorm"

	"cassette/pkg/storage"
	"cassette/pkg/storage/filesystem"
)

type Config struct {
	Repository repository.Repository
	Storage    storage.Storage
}

func FromEnvironment() (*Config, error) {
	r, err := gorm.NewSQLite("sessions/db.sqlite3")

	if err != nil {
		return nil, err
	}

	fs, err := filesystem.New("sessions")

	if err != nil {
		return nil, err
	}

	return &Config{
		Repository: r,
		Storage:    fs,
	}, nil
}
