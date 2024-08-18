package tarantstore

import "github.com/tarantool/go-tarantool"

type ValueStore struct {
	db *tarantool.Connection
}

func NewValueStore(db *tarantool.Connection) *ValueStore {
	return &ValueStore{
		db: db,
	}
}

// TODO:
func (s *ValueStore) SetKeys(keys map[interface{}]interface{}) error {
	panic("unimplemented")
}

// TODO:
func (s *ValueStore) GetKeys(keys map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	panic("unimplemented")
}
