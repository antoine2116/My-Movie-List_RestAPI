package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func NewValidationError(err error) error {
	var ve validator.ValidationErrors
	var formattedErrors string

	if errors.As(err, &ve) {
		for _, e := range ve {
			formattedErrors += fmt.Sprintf("%s %s \n", e.Field(), getErrorMsg(e))
		}
	}

	return errors.New(formattedErrors)
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

	return "Unknown error"
}
