package users

import (
	"apous-films-rest-api/common"
	"apous-films-rest-api/test"
	"apous-films-rest-api/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRegistration(t *testing.T) {
	asserts := assert.New(t)

	common.InitTestDB()

	testCases := []test.APITestCase{
		{
			Path:             "/auth/register",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com", "password": "superpass", "password_confirmation": "superpass"}}`,
			ExpectedStatus:   201,
			ExpectedResponse: `{"data":{"email":"steve@gmail.com","token":"([\w-]*\.[\w-]*\.[\w-]*)"}}`,
			Message:          "Valid data should return StatusCreated (201)",
		},
		{
			Path:             "/auth/register",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com", "password": "superpass", "password_confirmation": "wrongconfirmation"}}`,
			ExpectedStatus:   400,
			ExpectedResponse: `{"error":"passwords do not match"}`,
			Message:          "Wrong password confirmation should return StatusBadRequest (400)",
		},
		{
			Path:             "/auth/register",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com", "password": "superpass", "password_confirmation": "superpass"}}`,
			ExpectedStatus:   409,
			ExpectedResponse: `{"error":"User already exists with the same email"}`,
			Message:          "Same email as another user should return StatusConflict (409)",
		},
	}

	r := test.MockRouter()
	auth := r.Group("/auth")
	AddUserAuthentication(auth)

	for _, testCase := range testCases {
		test.Endpoint(asserts, r, testCase)
	}

	common.FreeTestDB()
}

func TestUserLogin(t *testing.T) {
	asserts := assert.New(t)

	common.InitTestDB()

	// Insert test user
	user := User{
		Email:        "steve@gmail.com",
		PasswordHash: utils.HashPassword("superpass"),
	}

	CreateUser(&user)

	testCases := []test.APITestCase{
		{
			Path:             "/auth/login",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com", "password": "superpass"}}`,
			ExpectedStatus:   200,
			ExpectedResponse: `{"data":{"email":"steve@gmail.com","token":"([\w-]*\.[\w-]*\.[\w-]*)"}}`,
			Message:          "Valid data should return StatusOK (200)",
		},
		{
			Path:             "/auth/login",
			Method:           "POST",
			Body:             `{"user": {"email": "wrongemail@gmail.com", "password": "superpass"}}`,
			ExpectedStatus:   400,
			ExpectedResponse: `{"error":"invalid email or password"}`,
			Message:          "Wrong email should return StatusBadRequest (400)",
		},
		{
			Path:             "/auth/login",
			Method:           "POST",
			Body:             `{"user": {"email": "steve@gmail.com", "password": "wrongpass"}}`,
			ExpectedStatus:   400,
			ExpectedResponse: `{"error":"invalid email or password"}`,
			Message:          "Wrong password should return StatusBadRequest (400)",
		},
	}

	r := test.MockRouter()
	auth := r.Group("/auth")
	AddUserAuthentication(auth)

	for _, testCase := range testCases {
		test.Endpoint(asserts, r, testCase)
	}

	common.FreeTestDB()
}
