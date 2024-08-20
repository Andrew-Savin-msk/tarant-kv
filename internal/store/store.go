package store

import (
	"github.com/Andrew-Savin-msk/tarant-kv/internal/domain/models"
)

type Bind struct {
	Key   string
	Value interface{}
}

type ValueStore interface {
	SetKeys(data map[string]interface{}) ([]Bind, error)
	GetKeys(keys []string) (map[string]interface{}, []string, error)
}
 
type UserStore interface {
	FindUser(login string) (*models.User, error)
}
