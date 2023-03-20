package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_extractToken(t *testing.T) {
	asserts := assert.New(t)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{}

	ctx.Request.Header = http.Header{"Authorization": {"Bearer TEST"}}
	token := extractToken(ctx)
	asserts.Equal("TEST", token)

	ctx.Request.Header = http.Header{"Authorization": {"wrongformat"}}
	token = extractToken(ctx)
	asserts.Empty(token)
}

func Test_GetCurrentUser(t *testing.T) {
	asserts := assert.New(t)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Set(userKey, CurrentUser{
		ID:    "test1",
		Email: "steve@gmail.com",
	})

	user := GetCurrentUser(ctx)

	asserts.Equal("test1", user.ID)
	asserts.Equal("steve@gmail.com", user.Email)
}
