package validations

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validate takes an object and validates its fields based on struct tags.
// It returns an error if validation fails, otherwise returns nil.
func Validate(object interface{}) error {
	validate := validator.New()
	if err := validate.Struct(object); err != nil {
		return fmt.Errorf(buildErrorMessage(err))
	}

	return nil
}

// buildErrorMessage constructs a detailed error message from validation errors.
// It returns a formatted string describing which fields failed validation and why.
func buildErrorMessage(err error) string {
	var errorMessage string

	validationErrors, validCast := err.(validator.ValidationErrors)
	if !validCast {
		return fmt.Sprintf("failed to validate with error: %v", err)
	}

	for i, err := range validationErrors {
		reason := err.ActualTag()
		separator := ""

		if err.Param() != "" {
			reason += " " + err.Param()
		}

		if i > 0 {
			separator = " | "
		}

		errorMessage += fmt.Sprintf("%sinvalid field: '%s'. Reason: %s.", separator, err.Field(), reason)
	}

	return errorMessage
}
