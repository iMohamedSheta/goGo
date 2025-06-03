package validate

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validatorInstance *validator.Validate
	once              sync.Once
)

func Validator() *validator.Validate {
	once.Do(func() {
		validatorInstance = validator.New()
	})
	return validatorInstance
}
