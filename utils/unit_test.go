package utils

import (
	"apous-films-rest-api/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Token
func TestGenerateJwt(t *testing.T) {
	asserts := assert.New(t)

	config.LoadConfiguration("../")
	token := GenerateJWT("id")

	asserts.Len(token, 143, "token length should be 143")
}

// Password
func TestHashPassword(t *testing.T) {
	asserts := assert.New(t)

	hash := HashPassword("mysuperpassword")

	asserts.Len(hash, 60, "hash length should be 60")
}

func TestCompareHashAndPassword(t *testing.T) {
	asserts := assert.New(t)

	password := "mysuperpassword"
	hash := HashPassword(password)

	asserts.NoError(CompareHashAndPassword(hash, "mysuperpassword"), "password should be validated")
	asserts.Error(CompareHashAndPassword(hash, "wrongpassword"), "password should not be valideted")
}
