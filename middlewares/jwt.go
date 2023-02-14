package middlewares

import (
	"apous-films-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := utils.VerifyToken(c); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Unauthorized : Invalid token"})
			return
		}

		c.Next()
	}
}
