package converter

import (
	"github.com/alfianyulianto/pds-service/internal/entity"
	"github.com/alfianyulianto/pds-service/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Email:           user.Email,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Password:        user.Password,
		Phone:           user.Phone,
		Avatar:          user.Avatar,
		IsActive:        user.IsActive,
		LastLoginAt:     user.LastLoginAt,
		Role:            user.Role,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		DeletedAt:       user.DeletedAt,
	}
}

func CreateRequestToUser(request *model.CreateUserRequest) *entity.User {
	return &entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Phone:    request.Phone,
		IsActive: request.IsActive,
	}
}

func UpdateRequestToUser(user *entity.User, request *model.UpdateUserRequest) *entity.User {
	user.Name = request.Name
	user.Email = request.Email
	user.Password = request.Password
	user.Phone = request.Phone
	user.IsActive = request.IsActive

	return user

}
