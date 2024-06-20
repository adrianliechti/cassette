package main

import (
	_ "embed"
	"net/http"

	"cassette/config"
	"cassette/pkg/server"
)

func main() {
	cfg, err := config.FromEnvironment()

	if err != nil {
		panic(err)
	}

	s := server.New(cfg)

	if err := http.ListenAndServe(":3000", s); err != nil {
		panic(err)
	}
}
