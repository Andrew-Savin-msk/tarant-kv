package apiserver

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (s *server) handleRegister() http.HandlerFunc {
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

		pHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrHashingPassword)
		}

		err = s.userStore.SaveUser(req.Username, pHash)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, ErrInvalidCredentials)
			return
		}

		s.respond(w, r, http.StatusOK, map[string]string{"status": "success"})
	})
}

func (s *server) handleLogin() http.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		panic("unimplemented endpoint")

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
