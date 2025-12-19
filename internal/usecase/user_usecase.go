package usecase

import (
	"context"
	"github.com/alfianyulianto/pds-service/internal/entity"
	"github.com/alfianyulianto/pds-service/internal/model"
	"github.com/alfianyulianto/pds-service/internal/model/converter"
	"github.com/alfianyulianto/pds-service/internal/repository"
	"github.com/alfianyulianto/pds-service/internal/utils"
	"github.com/alfianyulianto/pds-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Create(ctx context.Context, request *model.CreateUserRequest) (*model.UserResponse, error)
	Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error)
	List(ctx context.Context, request *model.SearchUserRequest) (*[]model.UserResponse, *response.Pagination, error)
	FindById(ctx context.Context, id any) (*model.UserResponse, error)
	Delete(ctx context.Context, id any) error
}

type userUseCase struct {
	*BaseUseCase
	UserRepository repository.UserRepository
}

func NewUserUseCase(baseUseCase *BaseUseCase, userRepository repository.UserRepository) UserUseCase {
	return &userUseCase{BaseUseCase: baseUseCase, UserRepository: userRepository}
}

func (u *userUseCase) Create(ctx context.Context, request *model.CreateUserRequest) (*model.UserResponse, error) {
	var success bool
	var filePath string

	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithField("action", "create user").WithError(err).Warn("Failed to validate request body")
		return nil, err
	}

	user := converter.CreateRequestToUser(request)

	if request.Avatar != nil {
		var err error
		filePath, err = u.Storage.UploadFile(request.Avatar, "user")
		if err != nil {
			u.Log.WithField("action", "create user").WithError(err).Error("Failed to upload avatar file")
			return nil, err
		}

		defer utils.CleanUpFilesOnFail(u.Storage, &success, filePath)
		user.Avatar = &filePath
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.WithField("action", "login").WithError(err).Error("Failed to hash password")
		return nil, fiber.ErrInternalServerError
	}

	user.Password = string(password)

	if err = u.UserRepository.Create(tx, user); err != nil {
		u.Log.WithField("action", "create user").WithError(err).Error("Failed to create user")
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		u.Log.WithField("action", "create user").WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	user.Avatar = utils.BuildFileURL(u.Config, filePath)

	success = true
	return converter.UserToResponse(user), nil
}

func (u *userUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	var success bool
	var filePath string

	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.ID); err != nil {
		u.Log.WithField("action", "update user").WithError(err).Error("Failed to find user")
		return nil, fiber.NewError(fiber.StatusNotFound, "User data not found")
	}

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithField("action", "update user").WithError(err).Warn("Failed to validate request body")
		return nil, err
	}

	user = converter.UpdateRequestToUser(user, request)

	if request.Avatar != nil {
		var err error
		filePath, err = u.Storage.UploadFile(request.Avatar, "user")
		if err != nil {
			u.Log.WithField("action", "update user").WithError(err).Error("Failed to upload avatar file")
			return nil, err
		}

		defer utils.CleanUpFilesOnFail(u.Storage, &success, filePath)

		user.Avatar = &filePath
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			u.Log.WithField("action", "login").WithError(err).Error("Failed to hash password")
			return nil, fiber.ErrInternalServerError
		}

		user.Password = string(password)
	}

	if err := u.UserRepository.Update(tx, user); err != nil {
		u.Log.WithField("action", "update user").WithError(err).Error("Failed to update user")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithField("action", "update user").WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	user.Avatar = utils.BuildFileURL(u.Config, filePath)

	success = true
	return converter.UserToResponse(user), nil
}

func (u *userUseCase) List(ctx context.Context, request *model.SearchUserRequest) (*[]model.UserResponse, *response.Pagination, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithField("action", "list user").WithError(err).Warn("Failed to validate request body")
		return nil, nil, err
	}

	users, total, err := u.UserRepository.FindAll(tx, request)
	if err != nil {
		u.Log.WithField("action", "list user").WithError(err).Error("Failed to find users")
		return nil, nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		u.Log.WithField("action", "list user").WithError(err).Error("Failed to commit transaction")
		return nil, nil, fiber.ErrInternalServerError
	}

	responses := make([]model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *converter.UserToResponse(&user)
	}

	paginationMeta := response.ToPaginated(request.Page, request.PageSize, total)
	return &responses, paginationMeta, nil
}

func (u *userUseCase) FindById(ctx context.Context, id any) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, id); err != nil {
		u.Log.WithField("action", "find user").WithError(err).Error("Failed to find user")
		return nil, fiber.NewError(fiber.StatusNotFound, "User data not found")
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithField("action", "find user").WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (u *userUseCase) Delete(ctx context.Context, id any) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, id); err != nil {
		u.Log.WithField("action", "delete user").WithError(err).Error("Failed to find user")
		return fiber.NewError(fiber.StatusNotFound, "User data not found")
	}

	if err := u.UserRepository.SoftDelete(tx, user); err != nil {
		u.Log.WithField("action", "delete user").WithError(err).Error("Failed to delete user")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithField("action", "delete user").WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	return nil
}
