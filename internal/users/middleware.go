package users

import (
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

type CurrentUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

const (
	userKey = "curr_user"
)

func GetCurrentUser(ctx *gin.Context) CurrentUser {
	return ctx.Value(userKey).(CurrentUser)
}

func JwtAuthentication(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Extract token from request
		stringToken := extractToken(ctx)

		if stringToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized : Missing token or invalid format"})
			return
		}

		// Validate token
		token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexptected signing method : %v", token.Header["alg"])
			}

			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized : Invalid token"})
			return
		}

		// Set current user
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx.Set(userKey, CurrentUser{
				ID:    claims["id"].(string),
				Email: claims["email"].(string),
			})
		}

		ctx.Next()
	}
}

func MockJWTAuthentication(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Header.Get("Authorization") != "TEST" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized : Invalid token"})
		}

		ctx.Set(userKey, CurrentUser{
			ID:    "test1",
			Email: "steve@gmail.com",
		})

		ctx.Next()
	}
}
