package repository

import (
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type Event any

type Session struct {
	ID string

	Created time.Time
}

type Repository struct {
	root string
}

func New() (*Repository, error) {
	root := "sessions"

	r := &Repository{
		root: root,
	}

	return r, nil
}

func (r *Repository) Sessions() ([]Session, error) {
	result := make([]Session, 0)

	return result, filepath.WalkDir(r.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()

		if err != nil {
			return nil
		}

		session := Session{
			ID: info.Name(),

			Created: info.ModTime(),
		}

		result = append(result, session)

		return nil
	})
}

func (r *Repository) Session(id string) (*Session, error) {
	path := filepath.Join(r.root, id)

	info, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	return &Session{
		ID: id,

		Created: info.ModTime(),
	}, nil
}

func (r *Repository) SessionEvents(id string) ([]Event, error) {
	path := filepath.Join(r.root, id)

	f, err := os.OpenFile(path, os.O_RDONLY, 0)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	var result []Event

	d := json.NewDecoder(f)

	for {
		var events []Event

		if err := d.Decode(&events); err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		result = append(result, events...)
	}

	return result, nil
}

func (r *Repository) AppendSessionEvents(id string, events ...Event) error {
	if len(events) == 0 {
		return nil
	}

	path := filepath.Join(r.root, id)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		return err
	}

	defer f.Close()

	if err := json.NewEncoder(f).Encode(events); err != nil {
		return err
	}

	return nil
}
