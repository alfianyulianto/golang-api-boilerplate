package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID              uuid.UUID      `gorm:"column:id;primaryKey"`
	Name            string         `gorm:"column:name;not null"`
	Email           string         `gorm:"column:email;not null"`
	EmailVerifiedAt *time.Time     `gorm:"column:email_verified_at"`
	Password        string         `gorm:"column:password;not null"`
	Phone           *string        `gorm:"column:phone"`
	Avatar          *string        `gorm:"column:avatar"`
	IsActive        bool           `gorm:"column:is_active"`
	LastLoginAt     *time.Time     `gorm:"column:last_login_at"`
	Role            string         `gorm:"column:role;default:User"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
