package test

import (
	"apous-films-rest-api/router"

	"github.com/gin-gonic/gin"
)

func MockRouter() *gin.Engine {
	return router.New()
}
