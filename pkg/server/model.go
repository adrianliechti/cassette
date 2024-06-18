package server

import "time"

type Session struct {
	ID string `json:"id"`

	Created time.Time `json:"created"`
}
