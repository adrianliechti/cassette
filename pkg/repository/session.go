package repository

import "time"

type Event any

type Session struct {
	ID string

	Created time.Time

	events []Event
}

func (s *Session) Events() []Event {
	if s.events == nil {
		return []Event{}
	}

	return s.events
}

func (s *Session) AppendEvents(events ...Event) error {
	s.events = append(s.events, events...)

	return nil
}
