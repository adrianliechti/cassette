package config

import (
	"cassette/pkg/storage"
	"cassette/pkg/storage/filesystem"
	"os"
	"path/filepath"

	"cassette/pkg/repository"
	"cassette/pkg/repository/gorm"
)

type Config struct {
	Storage    storage.Storage
	Repository repository.Repository
}

func FromEnvironment() (*Config, error) {
	path := os.Getenv("DATA_PATH")

	if path == "" {
		path = "sessions"
	}

	s, err := filesystem.New(path)

	if err != nil {
		return nil, err
	}

	r, err := gorm.NewSQLite(filepath.Join(path, "db.sqlite3"))

	if err != nil {
		return nil, err
	}

	return &Config{
		Storage:    s,
		Repository: r,
	}, nil
}
