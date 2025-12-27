package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using go-playground/validator
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// ParseAndValidate parses JSON body and validates the struct
func ParseAndValidate(c *fiber.Ctx, s interface{}) error {
	if err := c.BodyParser(s); err != nil {
		return err
	}
	return ValidateStruct(s)
}

// FormatValidationErrors formats validation errors into a readable string
func FormatValidationErrors(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errStr string
		for _, e := range validationErrors {
			if errStr != "" {
				errStr += "; "
			}
			errStr += formatFieldError(e)
		}
		return errStr
	}
	return err.Error()
}

// formatFieldError formats a single field error
func formatFieldError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return e.Field() + " is required"
	case "email":
		return e.Field() + " must be a valid email"
	case "min":
		return e.Field() + " must be at least " + e.Param() + " characters"
	case "max":
		return e.Field() + " must be at most " + e.Param() + " characters"
	case "eqfield":
		return e.Field() + " must match " + e.Param()
	default:
		return e.Field() + " failed validation: " + e.Tag()
	}
}
