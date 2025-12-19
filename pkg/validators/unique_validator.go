package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
	"strings"
)

func Unique(db *gorm.DB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		param := fl.Param()
		value := fl.Field().Interface()
		parent := fl.Parent()

		parts := strings.Split(param, ".")

		if len(parts) < 2 {
			return false
		}

		table := parts[0]
		column := parts[1]
		ignoreFields := parts[2:]

		if value == "" {
			return false
		}

		var count int64
		query := db.Table(table).Where(column+" = ?", value)

		for _, fieldName := range ignoreFields {
			f := parent.FieldByName(fieldName)
			if !f.IsValid() {
				return false
			}

			query = query.Where(strcase.ToSnake(fieldName)+" <> ?", f.Interface())
		}

		query = query.Count(&count)
		return count == 0
	}
}
