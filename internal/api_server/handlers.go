package apiserver

import (
	"encoding/json"
	"net/http"
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



		s.respond(w, r, http.StatusOK, map[string]string{"status": "success"})
	})
}

func (s *server) handleWriteKeys() http.HandlerFunc {
	type request struct {
		Data map[interface{}]interface{} `json:"data"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.value_store.SetKeys(req.Data)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]string{"status": "success"})
	})
}

func (s *server) handleReadKeys() http.HandlerFunc {
	type request struct {
		Data map[interface{}]interface{} `json:"data"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		res, err := s.value_store.GetKeys(req.Data)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInternalDbError)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]interface{}{"data": res})
	})
}
