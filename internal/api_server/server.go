package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Andrew-Savin-msk/tarant-kv/internal/store"
	"github.com/google/uuid"
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
	router *mux.Router
	logger *logrus.Logger
	store  store.Store
}

func newServer(st store.Store, logger *logrus.Logger) *server {
	srv := &server{
		router: mux.NewRouter(),
		logger: logger,
		store:  st,
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
	s.router.Use(s.authenticateUser)

}

// Sets unical id for every request
func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-Id", id)

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

// logRequest loggin any request and it's responce
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remout_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWrighter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed in %v %s in %v",
			rw.code,
			http.StatusText((rw.code)),
			time.Now().Sub(start),
		)
	})
}

// authenticateUser autentificates user by token in Authorisation header
func (s *server) authenticateUser(next http.Handler) http.Handler {
	// TODO:
	panic("unimplemented")
	// return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	session, err := s.sessionStore.Get(r, sessionName)
	// 	if err != nil {
	// 		s.error(w, r, http.StatusInternalServerError, err)
	// 		return
	// 	}

	// 	id, ok := session.Values["user_id"]
	// 	if !ok {
	// 		s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
	// 		return
	// 	}

	// 	u, err := s.store.User().Find(id.(int))
	// 	if err != nil {
	// 		s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
	// 		return
	// 	}

	// 	next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	// })
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
