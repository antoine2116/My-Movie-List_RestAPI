package routes

import (
	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine) {
	addTestsRoutes(r)
}
