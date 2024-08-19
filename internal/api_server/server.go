package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Andrew-Savin-msk/tarant-kv/internal/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "session-id"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not auntificated")
)

type ctxKey int8

type server struct {
	router     *mux.Router
	logger     *logrus.Logger
	valueStore store.ValueStore
	userStore  store.UserStore
	sessionStore sessions.CookieStore
	tokenTTL   time.Duration
}

func newServer(vSt store.ValueStore, uSt store.UserStore, logger *logrus.Logger, tokenTTL time.Duration) *server {
	srv := &server{
		router:     mux.NewRouter(),
		logger:     logger,
		valueStore: vSt,
		userStore:  uSt,
		sessionStore: sessions.NewStore()
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
