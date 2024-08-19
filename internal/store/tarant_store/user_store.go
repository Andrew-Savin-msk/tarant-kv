package tarantstore

import (
	"github.com/Andrew-Savin-msk/tarant-kv/internal/domain/models"
	"github.com/Andrew-Savin-msk/tarant-kv/internal/store"
	"github.com/tarantool/go-tarantool"
)

const schemaNameUsers = "users"

type UserStore struct {
	db *tarantool.Connection
}

func NewUserStore(db *tarantool.Connection) *UserStore {
	return &UserStore{
		db: db,
	}
}

// func (s *UserStore) SaveUser(login string, pHash []byte) error {
// 	panic("unimplemented")
// }

func (s *UserStore) FindUser(login string) (*models.User, error) {
	const op = "userstore.FindUser"

	resp, err := s.db.Select(schemaNameUsers, "primary", 0, 2, tarantool.IterEq, []interface{}{login})
	if err != nil {
		return nil, err
	}
	if len(resp.Tuples()) == 0 {
		return nil, store.ErrRecordNotFound
	}
	u := &models.User{
		Email:        resp.Tuples()[0][0].(string),
		PasswordHash: []byte(resp.Tuples()[0][1].(string)),
	}
	return u, nil
}

// func (s *UserStore) SaveToken(login, token string) error {
// 	panic("unimplemented")
// }
