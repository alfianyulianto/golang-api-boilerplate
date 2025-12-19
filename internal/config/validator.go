package config

import (
	validators2 "github.com/alfianyulianto/pds-service/pkg/validators"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func NewValidator(db *gorm.DB) *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("exists", validators2.Exists(db))
	validate.RegisterValidation("image", validators2.Image(db))
	validate.RegisterValidation("unique", validators2.Unique(db))
	validate.RegisterValidation("size", validators2.Size(db))
	validate.RegisterValidation("match_password", validators2.MatchPassword(db))

	return validate
}
