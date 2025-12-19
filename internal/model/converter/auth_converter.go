package converter

import (
	"github.com/alfianyulianto/pds-service/internal/entity"
	"github.com/alfianyulianto/pds-service/internal/model"
)

func RegisterRequestToUser(request *model.RegisterUserRequest) *entity.User {
	return &entity.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	}
}
