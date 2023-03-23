package users

import (
	"apous-films-rest-api/internal/test"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_handlers_Register_and_Login(t *testing.T) {
	asserts := assert.New(t)

	r := test.MockRouter()

	RegisterHandlers(r.Group(""),
		NewService(&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{}),
		"http://localhost:3000",
	)

	testCases := []test.APITestCase{
		{
			Path:             "/register",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com", "password": "pass", "passwordConfirmation": "pass"}}`,
			ExpectedStatus:   201,
			Message:          "Valid data should return StatusCreated (201)",
		},
		{
			Path:             "/register",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com", "password": "pass", "passwordConfirmation": "wrongconfirmation"}}`,
			ExpectedStatus:   400,
			Message:          "Wrong password confirmation should return StatusBadRequest (400)",
		},
		{
			Path:             "/register",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com", "password": "pass", "passwordConfirmation": "pass"}}`,
			ExpectedStatus:   409,
			Message:          "Same email as another user should return StatusConflict (409)",
		},
		{
			Path:             "/login",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com", "password": "pass"}}`,
			ExpectedStatus:   200,
			Message:          "Valid data should return StatusOK (200)",
		},
		{
			Path:             "/login",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com"}}`,
			ExpectedStatus:   400,
			Message:          "Invalid data (missing password) should return StatusBadRequest (400)",
		},
		{
			Path:             "/login",
			Method:           "POST",
			Body:             `{"user": {"email": "notsteve@gmail.com", "password": "pass"}}`,
			ExpectedStatus:   401,
			Message:          "Wrong email should return StatusUnauthorized (401)",
		},
		{
			Path:             "/login",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com", "password": "wrongpass"}}`,
			ExpectedStatus:   401,
			Message:          "Wrong password should return StatusUnauthorized (401)",
		},
	}

	for _, testCase := range testCases {
		test.Endpoint(asserts, r, testCase)
	}
}

func Test_handlers_GoogleLogin(t *testing.T) {
	asserts := assert.New(t)

	r := test.MockRouter()

	RegisterHandlers(r.Group(""),
		NewService(&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{}),
		"http://localhost:3000",
	)

	testCases := []test.APITestCase{
		{
			Path:           "/google/callback?code=authorization_code",
			Method:         "GET",
			Body:           "",
			ExpectedStatus: 308,
			Message:        "Valid code should return StatusPermanentRedirect (308)",
		},
		{
			Path:           "/google/callback",
			Method:         "GET",
			Body:           "",
			ExpectedStatus: 400,
			Message:        "Missing code should return StatusBadRequest (400)",
		},
		{
			Path:           "/google/callback?code=invalid_code",
			Method:         "GET",
			Body:           "",
			ExpectedStatus: 401,
			Message:        "Invalid code should return StatusUnauthorized (401)",
		},
	}

	for _, testCase := range testCases {
		test.Endpoint(asserts, r, testCase)
	}
}

func Test_handlers_GitHubLogin(t *testing.T) {
	asserts := assert.New(t)

	r := test.MockRouter()

	RegisterHandlers(r.Group(""),
		NewService(&mockRepository{}, "secret", 1000, mockProvider{}, mockProvider{}),
		"http://localhost:3000",
	)

	testCases := []test.APITestCase{
		{
			Path:           "/github/callback?code=authorization_code",
			Method:         "GET",
			Body:           "",
			ExpectedStatus: 308,
			Message:        "Valid code should return StatusPermanentRedirect (308)",
		},
		{
			Path:           "/github/callback",
			Method:         "GET",
			Body:           "",
			ExpectedStatus: 400,
			Message:        "Missing code should return StatusBadRequest (400)",
		},
		{
			Path:           "/github/callback?code=invalid_code",
			Method:         "GET",
			Body:           "",
			ExpectedStatus: 401,
			Message:        "Invalid code should return StatusUnauthorized (401)",
		},
	}

	for _, testCase := range testCases {
		test.Endpoint(asserts, r, testCase)
	}
}

func Test_handlers_Profile(t *testing.T) {
	asserts := assert.New(t)

	r := test.MockRouter()
	grp := r.Group("")
	grp.Use(MockJWTAuthentication("secret"))

	RegisterAuthenticatedHandlers(grp)

	testCases := []test.APITestCase{
		{
			Path:             "/profile",
			Method:           "GET",
			Header:           http.Header{"Authorization": {"TEST"}},
			ExpectedStatus:   200,
			Message:          "Valid token should return StatusOK (200) and the current user",
			ExpectedResponse: `{"id":"test1","email":"steve@gmail.com"}`,
		},
		{
			Path:           "/profile",
			Method:         "GET",
			Header:         http.Header{"Authorization": {"WRONG TOKEN"}},
			ExpectedStatus: 401,
			Message:        "Invalid token return StatusUnauthorized (401)",
		},
	}

	for _, testCase := range testCases {
		test.Endpoint(asserts, r, testCase)
	}

}
