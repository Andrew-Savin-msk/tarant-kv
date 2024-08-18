package store

import (
	"github.com/Andrew-Savin-msk/tarant-kv/internal/domain/models"
)

type ValueStore interface {
	SetKeys(keys map[interface{}]interface{}) error
	GetKeys(keys map[interface{}]interface{}) (map[interface{}]interface{}, error)
}

type UserStore interface {
	SaveUser(login string, pHash []byte) error
	FindUser(login string) (*models.User, error)
	SaveToken(login, token string) error
}
