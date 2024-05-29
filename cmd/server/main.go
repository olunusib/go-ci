package main

import (
	"log"

	"github.com/olunusib/go-ci/internal/config"
	"github.com/olunusib/go-ci/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	port := cfg.PORT
	log.Printf("Starting server on port %s", port)
	server.StartServer(cfg)
}
