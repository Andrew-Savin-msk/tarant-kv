package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Andrew-Savin-msk/tarant-kv/internal/store"
	"github.com/gorilla/mux"
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
	router      *mux.Router
	logger      *logrus.Logger
	value_store store.ValueStore
	user_store  store.UserStore
}

func newServer(st store.ValueStore, logger *logrus.Logger) *server {
	srv := &server{
		router:      mux.NewRouter(),
		logger:      logger,
		value_store: st,
	}

	srv.configureRouter()

	return srv
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.PathPrefix("/api")

	// TODO: Now it's off for testing
	// s.router.Use(s.authenticateUser)
	s.router.HandleFunc("/login", s.handleLogin()).Methods("POST")
	s.router.HandleFunc("/write", s.handleWriteKeys()).Methods("PUT")
	s.router.HandleFunc("/read", s.handleReadKeys()).Methods("GET")

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
