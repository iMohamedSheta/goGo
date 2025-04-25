package validate

import "github.com/go-playground/validator/v10"

var validatorInstance *validator.Validate

func Validator() *validator.Validate {
	if validatorInstance == nil {
		validatorInstance = validator.New()
	}
	return validatorInstance
}
