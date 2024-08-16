package tarantstore

import "github.com/tarantool/go-tarantool"

type Store struct {
	db *tarantool.Connection
}

func New(db *tarantool.Connection) *Store {
	return &Store{
		db: db,
	}
}

// TODO:
func (s *Store) SetKeys(keys map[interface{}]interface{}) error {
	panic("unimplemented")
}

// TODO:
func (s *Store) GetKeys(keys map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	panic("unimplemented")
}
