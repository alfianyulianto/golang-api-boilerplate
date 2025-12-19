package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/alfianyulianto/pds-service/internal/entity"
	"github.com/alfianyulianto/pds-service/internal/model"
)

type UserRepository interface {
	Create(db *gorm.DB, user *entity.User) error
	Update(db *gorm.DB, user *entity.User) error
	FindById(db *gorm.DB, user *entity.User, id any) error
	SoftDelete(db *gorm.DB, user *entity.User) error
	HardDelete(db *gorm.DB, user *entity.User) error
	FindAll(db *gorm.DB, request *model.SearchUserRequest) ([]entity.User, int64, error)
	FindByEmail(db *gorm.DB, user *entity.User, email string) error
}

type userRepository struct {
	Repository[entity.User]
	Log *logrus.Entry
}

func NewUserRepository(log *logrus.Entry) UserRepository {
	return &userRepository{Log: log}
}

func (r *userRepository) FindAll(db *gorm.DB, request *model.SearchUserRequest) ([]entity.User, int64, error) {
	var users []entity.User
	if err := db.Scopes(r.Filter(request)).Offset((request.Page - 1) * request.PageSize).Limit(request.PageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := db.Model(new(entity.User)).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

func (r *userRepository) Filter(request *model.SearchUserRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if request.Search != "" {
			search := "%" + request.Search + "%"

			tx = tx.Where("name like ?", "%"+search+"%").
				Or("email like ?", "%"+search+"%").
				Or("phone like ?", "%"+search+"%")
		}

		if request.IsActive != "" {
			tx = tx.Where("is_active = ?", request.IsActive)
		}

		if request.Role != "" {
			tx = tx.Where("role = ?", request.Role)
		}

		if request.OrderBy != "" && request.OrderDir != "" {
			tx = tx.Order(request.OrderBy + " " + request.OrderDir)
		} else {
			tx = tx.Order("created_at desc")
		}

		return tx
	}
}

func (r *userRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).Take(user).Error
}
