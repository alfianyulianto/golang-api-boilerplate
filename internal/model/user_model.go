package model

import (
	"github.com/alfianyulianto/pds-service/pkg/response"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"time"
)

type UserResponse struct {
	ID              uuid.UUID      `json:"id"`
	Name            string         `json:"name"`
	Email           string         `json:"email"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at"`
	Password        string         `json:"-"`
	Phone           *string        `json:"phone"`
	Avatar          *string        `json:"avatar"`
	IsActive        bool           `json:"is_active"`
	LastLoginAt     *time.Time     `json:"last_login_at"`
	Role            string         `json:"role"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at"`
}

type CreateUserRequest struct {
	Name            string                `json:"name" form:"name" validate:"required,max=255"`
	Email           string                `json:"email" form:"email" validate:"required,email,max=100,unique=users.email"` // unique=table.column
	Password        string                `json:"password" form:"password" validate:"required,min=8,max=100"`
	ConfirmPassword string                `json:"confirm_password" form:"confirm_password" validate:"required,eqfield=Password"`
	Phone           *string               `json:"phone" form:"phone" validate:"omitempty,max=20"`
	Avatar          *multipart.FileHeader `json:"avatar" form:"avatar" validate:"omitempty,image,size=2"`
	IsActive        bool                  `json:"is_active" form:"is_active" validate:"boolean"`
}

type UpdateUserRequest struct {
	ID              uuid.UUID             `json:"id" form:"id" validate:"required,uuid"`
	Name            string                `json:"name" form:"name" validate:"required,max=255"`
	Email           string                `json:"email" form:"email" validate:"required,email,max=100,unique=users.email.ID"` // unique=table.column.ignore1.ignore2
	Password        string                `json:"password" form:"password" validate:"omitempty,min=8,max=100"`
	ConfirmPassword string                `json:"confirm_password" form:"confirm_password" validate:"required_with,eqfield=Password"`
	Phone           *string               `json:"phone" form:"phone" validate:"omitempty,max=20"`
	Avatar          *multipart.FileHeader `json:"avatar" form:"avatar" validate:"omitempty,image,size=2"`
	IsActive        bool                  `json:"is_active" form:"is_active" validate:"boolean"`
}

type SearchUserRequest struct {
	Search   string `json:"search" form:"search" validate:"omitempty"`
	IsActive string `json:"is_active" form:"is_active" validate:"omitempty,boolean"`
	Role     string `json:"role" form:"role" validate:"omitempty,oneof=Admin User"`
	OrderBy  string `json:"order_by" validate:"omitempty"`
	OrderDir string `json:"order_dir" validate:"omitempty,oneof=asc desc"`
	response.PaginationRequest
}

type GetUserRequest struct {
	ID uuid.UUID `json:"id" form:"id" validate:"required,uuid"`
}
