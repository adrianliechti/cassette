package repository

import "time"

type Repository struct {
	sessions map[string]*Session

	last string
}

func New() (*Repository, error) {
	r := &Repository{
		sessions: make(map[string]*Session),
	}

	return r, nil
}

func (r *Repository) Sessions() ([]Session, error) {
	result := make([]Session, 0)

	for _, s := range r.sessions {
		result = append(result, *s)
	}

	return result, nil
}

func (r *Repository) Session(id string) (*Session, error) {
	if id == "default" {
		id = r.last
	}

	if session, ok := r.sessions[id]; ok {
		return session, nil
	}

	session := &Session{
		ID: id,

		Created: time.Now(),
	}

	r.last = id
	r.sessions[id] = session

	return session, nil
}
