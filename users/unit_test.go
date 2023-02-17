package users

import (
	"apous-films-rest-api/common"
	"apous-films-rest-api/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRegistration(t *testing.T) {
	asserts := assert.New(t)

	testCases := []test.APITestCase{
		{
			Path:           "/auth/register",
			Method:         "POST",
			Body:           `{"user": {"username": "steve", "email": "steve@gmail.com", "password": "superpass", "password_confirmation": "superpass"}}`,
			ExpectedStatus: 201,
		},
	}

	r := test.MockRouter()
	auth := r.Group("/auth")
	AddRoutes(auth)

	for _, testCase := range testCases {
		test.Endpoint(asserts, r, testCase)
	}
}
