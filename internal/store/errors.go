package store

import "errors"

var (
	ErrRecordNotFound      = errors.New("no such record")
	ErrRecordAlreadyExists = errors.New("record already exists")
	ErrStartingTransaction = errors.New("unable to start transaction")
	ErrConnCLosed          = errors.New("connection closed")
)
