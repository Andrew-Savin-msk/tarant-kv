package store

type Store interface {
	SetKeys(keys map[interface{}]interface{}) (map[interface{}]interface{}, error)
	GetKeys(keys map[interface{}]interface{}) (map[interface{}]interface{}, error)
}
