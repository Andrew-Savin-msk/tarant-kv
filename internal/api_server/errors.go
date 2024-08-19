package apiserver

import "errors"

var (
	ErrInternalDbError    = errors.New("valid ending of operation is unable")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrHashingPassword    = errors.New("unable to hash password")
	ErrPanicHanding       = errors.New("invalid data lead to panic")
	ErrMissingToken       = errors.New("missing authorization token")
	ErrExpiredToken       = errors.New("token epired")
	ErrNotAuntificated    = errors.New("user is not auntificated")
	ErrInvalidToken       = errors.New("invalid acess token")
)
