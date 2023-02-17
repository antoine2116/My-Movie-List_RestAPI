package common

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	return r
}
