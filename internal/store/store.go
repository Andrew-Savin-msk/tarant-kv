package store

import "github.com/Andrew-Savin-msk/tarant-kv/internal/domain/models"

type ValueStore interface {
	SetKeys(keys map[interface{}]interface{}) error
	GetKeys(keys map[interface{}]interface{}) (map[interface{}]interface{}, error)
}

type UserStore interface {
	SaveUser(email string, pHash []byte) error
	FindUser(email string) (*models.User, error)
}
