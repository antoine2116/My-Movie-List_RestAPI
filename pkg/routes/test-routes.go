package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func addTestsRoutes(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

}
