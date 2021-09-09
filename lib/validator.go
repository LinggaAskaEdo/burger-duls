package lib

import "github.com/go-playground/validator/v10"

type Validator struct {
	*validator.Validate
}

// NewValidator sets up validator
func NewValidator() Validator {
	return Validator{validator.New()}
}
