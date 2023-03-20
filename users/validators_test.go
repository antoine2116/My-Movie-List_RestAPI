package users

import (
	"apous-films-rest-api/test"
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var registerTests = []struct {
	data          string
	expectedError string
}{
	{
		`{"user": {"email": "steve@gmail.com", "password": "pass", "passwordConfirmation": "pass"}}`,
		"",
	},
	{
		`{"user": {"email": "steve@gmail.com", "password": "pass"}}`,
		"PasswordConfirmation is required",
	},
	{
		`{"user": {"email": "steve@gmail.com", "password": "pass", "passwordConfirmation": "wrongpass"}}`,
		"passwords do not match",
	},
}

func Test_RegisterValidatorBind(t *testing.T) {
	asserts := assert.New(t)

	ctx := test.MockGinContext()

	for _, tc := range registerTests {
		ctx.Request.Body = io.NopCloser(bytes.NewBufferString(tc.data))

		v := &RegisterValidator{}
		err := v.Bind(ctx)

		if tc.expectedError == "" {
			asserts.Nil(err)
		} else {
			asserts.ErrorContains(err, tc.expectedError)
		}
	}
}

var loginTests = []struct {
	data          string
	expectedError string
}{
	{
		`{"user": {"email": "steve@gmail.com", "password": "pass"}}`,
		"",
	},
	{
		`{"user": {"email": "steve@gmail.com"}}`,
		"Password is required",
	},
}

func Test_LoginValidatorBind(t *testing.T) {
	asserts := assert.New(t)

	ctx := test.MockGinContext()

	for _, tc := range loginTests {
		ctx.Request.Body = io.NopCloser(bytes.NewBufferString(tc.data))

		v := &LoginValidator{}
		err := v.Bind(ctx)

		if tc.expectedError == "" {
			asserts.Nil(err)
		} else {
			asserts.ErrorContains(err, tc.expectedError)
		}
	}
}
