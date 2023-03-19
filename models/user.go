package models

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Provider     string
}
