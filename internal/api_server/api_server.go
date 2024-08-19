package apiserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Andrew-Savin-msk/tarant-kv/internal/config"
	"github.com/Andrew-Savin-msk/tarant-kv/internal/store"
	tarantstore "github.com/Andrew-Savin-msk/tarant-kv/internal/store/tarant_store"
	"github.com/sirupsen/logrus"
	"github.com/tarantool/go-tarantool"
	"golang.org/x/crypto/bcrypt"
)

// Start init's all connections and starts api's work
func Start(cfg *config.Config) error {
	// Get logger
	log := setLog(cfg.Srv.LogLevel)

	// Get value store
	vSt, err := connValueStore(&cfg.VDb)
	if err != nil {
		log.Fatalf("unable to connect value store (%s), ended with error: %s", cfg.VDb.Host, err)
	}

	// Get user store
	uSt, err := connUserStore(&cfg.UDb)
	if err != nil {
		log.Fatalf("unable to connect user store (%s), ended with error: %s", cfg.UDb.Host, err)
	}

	// Get server
	srv := newServer(vSt, uSt, log, cfg.Srv.TokenTTL)

	log.Infof("api strted work on port: %s", cfg.Srv.Port[1:])
	// Start listner
	err = http.ListenAndServe(cfg.Srv.Port, srv)

	if err != nil {
		log.Infof("api ended work with error: %w", err)
	} else {
		log.Info("api ended work")
	}

	return nil
}

// setLog set logger level by setted in config
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
	fmt.Printf("logger set in level: %s\n", level)
	return log
}

// connValueStore connects to a value store and returns store structure
func connValueStore(cfg *config.ValueDatabase) (store.ValueStore, error) {
	opts := tarantool.Opts{}

	con, err := tarantool.Connect(cfg.Host+cfg.Port, opts)
	if err != nil {
		return nil, err
	}

	_, err = con.Ping()
	if err != nil {
		return nil, err
	}

	return tarantstore.NewValueStore(con), nil
}

// connUserStore connects to a user store and returns store structure
func connUserStore(cfg *config.UserDatabase) (store.UserStore, error) {
	opts := tarantool.Opts{}

	con, err := tarantool.Connect(cfg.Host+cfg.Port, opts)
	if err != nil {
		return nil, err
	}

	_, err = con.Ping()
	if err != nil {
		return nil, err
	}

	pHash, err := bcrypt.GenerateFromPassword([]byte(cfg.DefPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = con.Replace("users", []interface{}{cfg.DefUser, string(pHash)})
	if err != nil {
		return nil, err
	}

	return tarantstore.NewUserStore(con), nil
}
