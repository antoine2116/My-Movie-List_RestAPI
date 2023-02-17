package test

import (
	"apous-films-rest-api/common"

	"github.com/gin-gonic/gin"
)

func MockRouter() *gin.Engine {
	return common.NewRouter()
}
