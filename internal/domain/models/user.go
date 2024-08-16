package models

type User struct {
	Email        string
	PasswordHash []byte
}
