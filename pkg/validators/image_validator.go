package validators

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
)

func Image(db *gorm.DB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		file, ok := fl.Field().Interface().(multipart.FileHeader)
		if !ok {
			return true
		}

		contentType := file.Header.Get("Content-Type")
		fmt.Println(contentType)
		return strings.HasPrefix(contentType, "image/")
	}
}
