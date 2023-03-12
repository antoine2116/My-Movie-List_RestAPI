package utils

import "github.com/gin-gonic/gin"

func SetCookieToken(c *gin.Context, token string) {
	c.SetCookie("token", token, 60*60*24, "/", "", true, false)
}
