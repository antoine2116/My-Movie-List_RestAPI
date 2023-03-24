package errors

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func Test_errors_getErrorMsg(t *testing.T) {
	asserts := assert.New(t)

	// Required
	msg := getErrorMsg(mockFieldError{tag: "required"})
	asserts.Equal("is required", msg)

	// Less than or equal to
	msg = getErrorMsg(mockFieldError{tag: "lte", param: "10"})
	asserts.Equal("should be less than or equal to 10", msg)

	// Greater than or equal to
	msg = getErrorMsg(mockFieldError{tag: "gte", param: "5"})
	asserts.Equal("should be greater than or equal to 5", msg)

	// Unknown
	msg = getErrorMsg(mockFieldError{tag: "unknown"})
	asserts.Equal("unknown error", msg)
}

func Test_errors_NewValidationError(t *testing.T) {
	asserts := assert.New(t)

	err := NewValidationError(errors.New("some error message"))
	asserts.NotNil(err)
	asserts.Equal("some error message", err.Message)
	asserts.Equal(0, len(err.Errors))

	err = NewValidationError(validator.ValidationErrors{
		mockFieldError{field: "name", tag: "required"},
		mockFieldError{field: "age", tag: "lte", param: "200"},
	})

	asserts.NotNil(err)
	asserts.Equal("name", err.Errors[0].Field)
	asserts.Equal("name is required", err.Errors[0].Error)
	asserts.Equal("age", err.Errors[1].Field)
	asserts.Equal("age should be less than or equal to 200", err.Errors[1].Error)

	err = NewValidationError(validator.ValidationErrors{
		mockFieldError{field: "name", tag: "unknown"},
	})

	asserts.NotNil(err)
	asserts.Equal("name", err.Errors[0].Field)
	asserts.Equal("name unknown error", err.Errors[0].Error)
}

func Test_errors_NewCommonError(t *testing.T) {
	asserts := assert.New(t)
	
	err := NewCommonError(errors.New("error message"))
	asserts.Equal("error message", err.Message)
	asserts.Len(err.Errors, 0)
}

type mockFieldError struct {
	validator.FieldError
	tag   string
	field string
	param string
}

func (e mockFieldError) Tag() string {
	return e.tag
}

func (e mockFieldError) Field() string {
	return e.field
}

func (e mockFieldError) Param() string {
	return e.param
}
