package apiserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Andrew-Savin-msk/tarant-kv/internal/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "Authorization"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

type server struct {
	router     *mux.Router
	logger     *logrus.Logger
	valueStore store.ValueStore
	userStore  store.UserStore
	sessionKey string
	tokenTTL   time.Duration
}

func newServer(vSt store.ValueStore, uSt store.UserStore, logger *logrus.Logger, sessionKey string, tokenTTL time.Duration) *server {
	srv := &server{
		router:     mux.NewRouter(),
		logger:     logger,
		valueStore: vSt,
		userStore:  uSt,
		sessionKey: sessionKey,
		tokenTTL:   tokenTTL,
	}

	srv.configureRouter()

	return srv
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Handle("/api/login", s.basePaths(s.handleLogin()))
	s.router.Handle("/api/write", s.protectedPaths(s.handleWriteKeys()))
	s.router.Handle("/api/read", s.protectedPaths(s.handleReadKeys()))
}

// Func for making call of respond func with Error pattern
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

// Universal func for sending any type of respond (Error, Responde, etc.)
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
