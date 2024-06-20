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

	Username string
	Password string
}

func FromEnvironment() (*Config, error) {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	if username == "" {
		username = "admin"
	}

	if password == "" {
		password = "admin"
	}

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

		Username: username,
		Password: password,
	}, nil
}
