package filesystem

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"cassette/pkg/storage"
)

var _ storage.Storage = &FileSystem{}

type FileSystem struct {
	root string
}

func New(root string) (*FileSystem, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}

	return &FileSystem{
		root: root,
	}, nil
}

func (fs *FileSystem) Events(session string) ([]storage.Event, error) {
	path := filepath.Join(fs.root, session)

	f, err := os.OpenFile(path, os.O_RDONLY, 0)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	var result []storage.Event

	d := json.NewDecoder(f)

	for {
		var events []storage.Event

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

func (fs *FileSystem) AppendEvents(session string, events ...storage.Event) error {
	if len(events) == 0 {
		return nil
	}

	path := filepath.Join(fs.root, session)

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

func (fs *FileSystem) DeleteDelete(session string) error {
	path := filepath.Join(fs.root, session)

	os.Remove(path)
	return nil
}
