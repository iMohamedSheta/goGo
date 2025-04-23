package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Validate(data map[string]interface{}, rules map[string]string, messages map[string]string) (bool, map[string]string) {
	validate := validator.New()
	errors := make(map[string]string)

	for field, rule := range rules {
		err := validate.Var(data[field], rule)
		if err != nil {
			// Get the first error for the field
			if validationErrs, ok := err.(validator.ValidationErrors); ok && len(validationErrs) > 0 {
				// Get the custom error message if available
				customMessage := fmt.Sprintf("Field %s failed on '%s'", field, validationErrs[0].Tag())

				// Check if custom message exists in the messages map
				if msg, exists := messages[field]; exists {
					customMessage = msg
				}

				errors[field] = customMessage
			} else {
				errors[field] = "Validation failed"
			}
		}
	}

	return len(errors) == 0, errors
}
