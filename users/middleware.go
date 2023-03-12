package users

import (
	"apous-films-rest-api/config"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func extractToken(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")

	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func updateUserContext(c *gin.Context, userId string) {
	var user User

	if err := FindUserById(&user, userId); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	c.Set("user", user)
}

func JwtAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		stringToken := extractToken(c)

		if stringToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized : Missing token or invalid format"})
			return
		}

		token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexptected signing method : %v", token.Header["alg"])
			}

			return []byte(config.Config.Server.Secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Unauthorized : Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := claims["sub"].(string)
			updateUserContext(c, userId)
		}

		c.Next()
	}
}
