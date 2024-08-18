package store

import "fmt"

var (
	ErrRecordNotFound      = fmt.Errorf("no such record")
	ErrRecordAlreadyExists = fmt.Errorf("record already exists")
)
