package validate

import (
	"fmt"
	"imohamedsheta/gocrud/pkg/contracts"

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

func ValidateRequest(r contracts.Validatable) (bool, map[string]string) {
	validate := validator.New()

	// Validate the struct
	if err := validate.Struct(r); err != nil {
		errors := make(map[string]string)
		messages := r.Messages()

		for _, valErr := range err.(validator.ValidationErrors) {
			key := valErr.Field() + "." + valErr.Tag()
			if msg, ok := messages[key]; ok {
				errors[valErr.Field()] = msg
			} else {
				errors[valErr.Field()] = valErr.Error()
			}
		}
		return false, errors
	}

	return true, nil
}
