package errors

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type CommonError struct {
	Message string       `json:"message"`
	Errors  []InputError `json:"errors"`
}

type InputError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func NewCommonError(err error) CommonError {
	return CommonError{Message: err.Error()}
}

func NewValidationError(err error) *CommonError {
	var ve validator.ValidationErrors

	// Simple error
	if !errors.As(err, &ve) {
		return &CommonError{Message: err.Error()}
	}

	// Building input errors from validation errors
	var iptErrors []InputError

	for _, e := range ve {
		iptErrors = append(iptErrors, InputError{
			Field: e.Field(),
			Error: fmt.Sprintf("%s %s", e.Field(), getErrorMsg(e)),
		})
	}

	return &CommonError{Errors: iptErrors}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "lte":
		return fmt.Sprintf("should be less than or equal to %s", fe.Param())
	case "gte":
		return fmt.Sprintf("should be greater than or equal to %s", fe.Param())
	}

	return "unknown error"
}
