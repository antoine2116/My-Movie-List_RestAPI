package utils

import (
	"apous-films-rest-api/config"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	tokenDuration = 15000
)

func GenerateJWT(userId string) string {
	conf := config.GetConfig()
	secret := []byte(conf.Server.Secret)

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * tokenDuration).Unix(),
		"iat": time.Now().Unix(),
		"sub": userId,
	}

	tokenString, err := token.SignedString(secret)

	if err != nil {
		panic(err)
	}

	return tokenString
}

func VerifyToken(c *gin.Context) error {
	conf := config.GetConfig()

	stringToken := extractToken(c)

	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexptected signing method : %v", token.Header["alg"])
		}

		return []byte(conf.Server.Secret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid Token")
	}

	return nil
}

func extractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
