package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

func ParseErrors(err error) map[string]ErrorMessage {
	errors := make(map[string]ErrorMessage)

	for _, validationError := range err.(validator.ValidationErrors) {
		errors[strcase.ToSnake(validationError.Field())] = ErrorMessage{
			Field:   strcase.ToSnake(validationError.Field()),
			Message: GetCustomeMessage(validationError),
			Tag:     validationError.Tag(),
			Value:   validationError.Value(),
		}
	}
	return errors
}
