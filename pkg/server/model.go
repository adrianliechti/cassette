package server

import "time"

type Event any

type Session struct {
	ID string `json:"id"`

	Created time.Time `json:"created"`
}
