package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash)
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
