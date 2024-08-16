package main

import (
	apiserver "github.com/Andrew-Savin-msk/tarant-kv/internal/api_server"
	"github.com/Andrew-Savin-msk/tarant-kv/internal/config"
)

func main() {
	cfg := config.Load()

	apiserver.Start(cfg)
}
