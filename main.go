package main

import (
	_ "embed"
	"net/http"

	"cassette/pkg/repository"
	"cassette/pkg/server"
)

func main() {
	r, err := repository.New()

	if err != nil {
		panic(err)
	}

	s := server.New(r)

	if err := http.ListenAndServe(":8080", s); err != nil {
		panic(err)
	}
}
