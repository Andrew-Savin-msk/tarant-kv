package apiserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Andrew-Savin-msk/tarant-kv/internal/config"
	"github.com/sirupsen/logrus"
)

func Start(cfg *config.Config) error {
	// Get logger

	log := setLog(cfg.Srv.LogLevel)

	// Get value store

	// Get user store

	// Get server
	srv := newServer(nil, log)

	// Start listner
	http.ListenAndServe(cfg.Srv.Port, srv)

	return nil
}

func setLog(level string) *logrus.Logger {
	log := logrus.New()
	switch strings.ToLower(level) {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	}
	fmt.Printf("logger set in level: %s \n", level)
	return log
}
