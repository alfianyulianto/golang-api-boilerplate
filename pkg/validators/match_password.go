package validators

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func MatchPassword(db *gorm.DB) validator.Func {
	return func(fl validator.FieldLevel) bool {
		param := fl.Param()
		value := fl.Field().String()
		parent := fl.Parent().FieldByName("ID")

		if param == "" {
			return false
		}

		if value == "" {
			return false
		}

		user := make(map[string]interface{})
		db.Table(param).Select("id, password").Where("id = ?", parent.Interface()).Take(&user)
		if len(user) == 0 {
			return false
		}

		err := bcrypt.CompareHashAndPassword([]byte(user["password"].(string)), []byte(value))
		if err != nil {
			return false
		}

		return true
	}
}
