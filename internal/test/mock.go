package test

import (
	"apous-films-rest-api/internal/router"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func MockRouter() *gin.Engine {
	return router.New()
}

func MockGinContext() *gin.Context {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{}

	return ctx
}
