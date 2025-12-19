package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   any    `json:"value"`
}

func GetCustomeMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "boolean":
		return strcase.ToSnake(fieldError.Field()) + " field must be true or false."
	case "datetime":
		return strcase.ToSnake(fieldError.Field()) + " field must match the format " + fieldError.Param() + "."
	case "eqfield":
		return strcase.ToSnake(fieldError.Field()) + " field must be equal to " + strcase.ToSnake(fieldError.Param()) + " field."
	case "exists":
		return strcase.ToSnake(fieldError.Field()) + " is invalid."
	case "image":
		return strcase.ToSnake(fieldError.Field()) + " field must be an image."
	case "match_password":
		return strcase.ToSnake(fieldError.Field()) + " is incorrect."
	case "max":
		return strcase.ToSnake(fieldError.Field()) + " field must not be greater than " + fieldError.Param() + " characters."
	case "min":
		return strcase.ToSnake(fieldError.Field()) + " field must be at least " + fieldError.Param() + " characters."
	case "number":
		return strcase.ToSnake(fieldError.Field()) + " field must be a number"
	case "numeric":
		return strcase.ToSnake(fieldError.Field()) + " field must be a number"
	case "oneof":
		return strcase.ToSnake(fieldError.Field()) + " is invalid."
	case "required":
		return strcase.ToSnake(fieldError.Field()) + " is required."
	case "size":
		return strcase.ToSnake(fieldError.Field()) + " field must not be greater than " + fieldError.Param() + " megabyte."
	case "unique", "unique_with_ignore":
		return strcase.ToSnake(fieldError.Field()) + " has already been taken."
	default:
		return fieldError.Error()
	}
}
