package apiserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Andrew-Savin-msk/tarant-kv/internal/lib/jwt"
	"github.com/Andrew-Savin-msk/tarant-kv/internal/store"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// setRequestID sets UUID for every request
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
		rw := &responseWriter{w, http.StatusOK}

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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(sessionName)
		if token == "" {
			s.error(w, r, http.StatusUnauthorized, ErrMissingToken)
			return
		}

		cl, err := jwt.DecodeJWT(token, s.sessionKey)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, ErrInvalidToken)
			return
		}

		u, err := s.userStore.FindUser(cl.Username)
		if u == nil || err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(w, r, http.StatusBadRequest, ErrNotAuntificated)
			}
			s.error(w, r, http.StatusBadRequest, ErrInvalidCredentials)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, cl.Username)))
	})
}

// recoverPanic panic recovering middleware
func (s *server) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger := s.logger.WithFields(logrus.Fields{
					"remout_addr": r.RemoteAddr,
					"request_id":  r.Context().Value(ctxKeyRequestID),
					"method":      r.Method,
					"URI":         r.RequestURI,
				})

				logger.Errorf("ended hadling by panic with error: %s", err)

				s.error(w, r, http.StatusInternalServerError, ErrPanicHanding)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// basePaths middleware wrapper for base paths
func (s *server) basePaths(next http.Handler) http.Handler {
	return s.setRequestID(s.recoverPanic(s.logRequest(next)))
}

// protectedPaths middleware wrapper for authorisation required paths
func (s *server) protectedPaths(next http.Handler) http.Handler {
	return s.authenticateUser(s.basePaths(next))
}
