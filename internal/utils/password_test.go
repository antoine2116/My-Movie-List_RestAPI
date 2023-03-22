package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_password_HashPassword(t *testing.T) {
	asserts := assert.New(t)

	hash := HashPassword("pass")
	asserts.Len(hash, 60, "hash length should be 60")
}

func TestCompareHashAndPassword(t *testing.T) {
	asserts := assert.New(t)

	password := "pass"
	hash := HashPassword(password)

	err := CompareHashAndPassword(hash, "pass")
	asserts.Nil(err)

	err = CompareHashAndPassword(hash, "wrongpass")
	asserts.Error(err)
}
