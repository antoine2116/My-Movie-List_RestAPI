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
		return "should be less than " + fe.Param()
	case "gte":
		return "should be greater than " + fe.Param()
	}

	return "Unknown error"
}
