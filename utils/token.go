package utils

import (
	"apous-films-rest-api/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(userId string) string {
	secret := []byte(config.Config.Server.Secret)

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * time.Duration(config.Config.Server.TokenDuration)).Unix(),
		"iat": time.Now().Unix(),
		"sub": userId,
	}

	tokenString, err := token.SignedString(secret)

	if err != nil {
		panic(err)
	}

	return tokenString
}
