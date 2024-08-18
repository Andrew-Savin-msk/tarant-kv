package tarantstore

import (
	"github.com/Andrew-Savin-msk/tarant-kv/internal/domain/models"
	"github.com/tarantool/go-tarantool"
)

type UserStore struct {
	db *tarantool.Connection
}

func NewUserStore(db *tarantool.Connection) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) SaveUser(login string, pHash []byte) error {
	panic("unimplemented")
}

func (s *UserStore) FindUser(login string) (*models.User, error) {
	panic("unimplemented")
}

func (s *UserStore) SaveToken(login, token string) error {
	panic("unimplemented")
}
