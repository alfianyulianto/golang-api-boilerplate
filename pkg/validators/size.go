package validators

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"mime/multipart"
	"strconv"
)

func Size(db *gorm.DB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		param, _ := strconv.Atoi(fl.Param())

		file, ok := fl.Field().Interface().(multipart.FileHeader)
		if !ok {
			return true
		}

		return file.Size < int64(param)*1024*1024
	}
}
