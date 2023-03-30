package users

import (
	"apous-films-rest-api/internal/errors"
	"apous-films-rest-api/internal/test"
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var registerTests = []struct {
	data          string
	expectedError *errors.CommonError
}{
	{
		`{"user": {"email": "steve@gmail.com", "password": "pass", "passwordConfirmation": "pass"}}`,
		nil,
	},
	{
		`{"user": {"email": "steve@gmail.com", "password": "pass"}}`,
		&errors.CommonError{
			Message: "",
			Errors: []errors.InputError{
				{Field: "PasswordConfirmation", Error: "PasswordConfirmation is required"},
			},
		},
	},
	{
		`{"user": {"email": "steve@gmail.com", "password": "pass", "passwordConfirmation": "wrongpass"}}`,
		&errors.CommonError{
			Message: ErrMismatchedPasswords.Error(),
			Errors:  nil,
		},
	},
}

func Test_validator_RegisterBind(t *testing.T) {
	asserts := assert.New(t)

	ctx := test.MockGinContext()

	for _, tc := range registerTests {
		ctx.Request.Body = io.NopCloser(bytes.NewBufferString(tc.data))

		v := &RegisterValidator{}
		err := v.Bind(ctx)

		if tc.expectedError == nil {
			asserts.Nil(err)
			continue
		}

		comErr := errors.NewValidationError(err)

		if tc.expectedError.Message != "" {
			asserts.Equal(tc.expectedError.Message, comErr.Message)
		}

		if tc.expectedError.Errors != nil {
			asserts.True(reflect.DeepEqual(tc.expectedError.Errors, comErr.Errors))
		}
	}
}

var loginTests = []struct {
	data          string
	expectedError *errors.CommonError
}{
	{
		`{"user": {"email": "steve@gmail.com", "password": "pass"}}`,
		nil,
	},
	{
		`{"user": {"email": "steve@gmail.com"}}`,
		&errors.CommonError{
			Message: "",
			Errors: []errors.InputError{
				{Field: "Password", Error: "Password is required"},
			},
		},
	},
}

func Test_validator_LoginBind(t *testing.T) {
	asserts := assert.New(t)

	ctx := test.MockGinContext()

	for _, tc := range loginTests {
		ctx.Request.Body = io.NopCloser(bytes.NewBufferString(tc.data))

		v := &LoginValidator{}
		err := v.Bind(ctx)

		if tc.expectedError == nil {
			asserts.Nil(err)
			continue
		}

		comErr := errors.NewValidationError(err)

		if tc.expectedError.Message != "" {
			asserts.Equal(tc.expectedError.Message, comErr.Message)
		}

		if tc.expectedError.Errors != nil {
			asserts.True(reflect.DeepEqual(tc.expectedError.Errors, comErr.Errors))
		}
	}
}
