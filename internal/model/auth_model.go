package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterUserRequest struct {
	Name            string `json:"name" form:"name" validate:"required,max=255"`
	Email           string `json:"email" form:"email" validate:"required,email,max=100,unique=users.email"` // unique=table.column
	Password        string `json:"password" form:"password" validate:"required,min=8,max=100"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginUserRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

type VerifyUserRequest struct {
	Token string `validate:"required"`
}

type UserClaimToken struct {
	ID   uuid.UUID
	Role string
	Type string
	jwt.RegisteredClaims
}

type UpdatePasswordRequest struct {
	ID                 uuid.UUID `validate:"required,exists=users.id"`
	OldPassword        string    `json:"old_password" form:"old_password" validate:"required,match_password=users"`
	NewPassword        string    `json:"new_password" form:"new_password" validate:"required,min=8,max=100,nefield"`
	ConfirmNewPassword string    `json:"confirm_new_password" form:"confirm_new_password" validate:"required,eqfield=NewPassword"`
}

type RequestPasswordResetRequest struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token              string `json:"token" form:"token" validate:"required"`
	NewPassword        string `json:"new_password" form:"new_password" validate:"required,min=8,max=100"`
	ConfirmNewPassword string `json:"confirm_new_password" form:"confirm_new_password" validate:"required,eqfield=NewPassword"`
}
