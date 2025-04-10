package common

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func Validate(data interface{}) []APIError {
	validationErrors := []APIError{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem APIError

			elem.Field = err.Field()
			elem.Message = err.Error()

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
