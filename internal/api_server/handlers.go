package apiserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Andrew-Savin-msk/tarant-kv/internal/store"
	"github.com/google/uuid"
)

func (s *server) handleLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if req.Username == "" || req.Password == "" {
			s.error(w, r, http.StatusBadRequest, ErrInvalidCredentials)
			return
		}

		u, err := s.userStore.FindUser(req.Username)
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				s.error(w, r, http.StatusBadRequest, ErrInvalidCredentials)
				return
			}
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// TODO: Generate token
		token := uuid.NewString()
		// TODO: Save token

		s.respond(w, r, http.StatusOK, u)
	})
}

func (s *server) handleWriteKeys() http.HandlerFunc {
	type request struct {
		Data map[string]interface{} `json:"data"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		uninserted, err := s.valueStore.SetKeys(req.Data)
		if err != nil {
			// TODO: needs alias
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if len(uninserted) != 0 {
			s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "partially_inserded", "not_inserted": uninserted})
		} else {
			s.respond(w, r, http.StatusOK, map[string]interface{}{"status": "success"})
		}
	})
}

func (s *server) handleReadKeys() http.HandlerFunc {
	type request struct {
		Keys []string `json:"keys"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		res, unfound, err := s.valueStore.GetKeys(req.Keys)
		if err != nil {
			// TODO: needs alias
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if len(unfound) != 0 {
			s.respond(w, r, http.StatusOK, map[string]interface{}{"data": res, "not_found": unfound})
		} else {
			s.respond(w, r, http.StatusOK, map[string]interface{}{"data": res})
		}
	})
}
