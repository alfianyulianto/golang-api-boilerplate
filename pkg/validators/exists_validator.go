package validators

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"strings"
)

func Exists(db *gorm.DB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		param := fl.Param()
		value := fl.Field().Interface()

		if value == "" {
			return false
		}

		parts := strings.Split(param, ".")

		if len(parts) != 2 {
			return false
		}

		table := parts[0]
		column := parts[1]

		var count int64
		db.Table(table).Where(column+" = ?", value).Count(&count)
		return count > 0
	}
}
