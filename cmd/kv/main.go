package main

import (
	"log"

	apiserver "github.com/Andrew-Savin-msk/tarant-kv/internal/api_server"
	"github.com/Andrew-Savin-msk/tarant-kv/internal/config"
)

func main() {
	cfg := config.Load()

	err := apiserver.Start(cfg)
	if err != nil {
		log.Fatalf("unable to start api? ended with error: %w", err)
	}
}
