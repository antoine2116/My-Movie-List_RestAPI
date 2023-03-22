package utils

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func Test_formats_getErrorMsg(t *testing.T) {
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
	asserts.Equal("Unknown error", msg)
}

func Test_formats_NewValidationError(t *testing.T) {
	asserts := assert.New(t)

	err := NewValidationError(validator.ValidationErrors{
		mockFieldError{field: "name", tag: "required"},
		mockFieldError{field: "age", tag: "lte", param: "200"},
	})

	asserts.NotNil(err)
	asserts.ErrorContains(err, "name is required \nage should be less than or equal to 200 \n")

	err = NewValidationError(validator.ValidationErrors{
		mockFieldError{field: "name", tag: "unknown"},
	})

	asserts.NotNil(err)
	asserts.ErrorContains(err, "name Unknown error \n")
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
