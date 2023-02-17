package utils

import (
	"apous-films-rest-api/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	tokenDuration = 24
)

func GenerateJWT(userId string) string {
	conf := config.LoadConfiguration("../")

	secret := []byte(conf.Server.Secret)

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * tokenDuration).Unix(),
		"iat": time.Now().Unix(),
		"sub": userId,
	}

	tokenString, err := token.SignedString(secret)

	if err != nil {
		panic(err)
	}

	return tokenString
}
