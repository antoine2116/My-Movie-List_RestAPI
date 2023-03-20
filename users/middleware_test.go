package users

import (
	"apous-films-rest-api/test"
	"fmt"
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

func Test_JwtAthentication(t *testing.T) {
	asserts := assert.New(t)

	r := test.MockRouter()
	grp := r.Group("")
	grp.Use(JwtAuthentication("secret"))
	grp.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	service := service{&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{}}
	validToken := service.generateJWT("test1", "steve@gmail.com")

	testCases := []test.APITestCase{
		{
			Path:           "/protected",
			Method:         "GET",
			Header:         http.Header{"Authorization": {fmt.Sprintf("Bearer %s", validToken)}},
			ExpectedStatus: 200,
			Message:        "Valid token should return StatusOK (200)",
		},
		{
			Path:           "/protected",
			Method:         "GET",
			ExpectedStatus: 401,
			Message:        "Missing Authorization header should return StatusUnauthorized (401)",
		},
		{
			Path:           "/protected",
			Method:         "GET",
			Header:         http.Header{"Authorization": {"invalid_token"}},
			ExpectedStatus: 401,
			Message:        "Invalid token should return StatusUnauthorized (401)",
		},
	}

	for _, testCase := range testCases {
		test.Endpoint(asserts, r, testCase)
	}
}
