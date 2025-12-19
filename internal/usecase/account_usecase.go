package usecase

import (
	"context"
	"github.com/alfianyulianto/pds-service/internal/entity"
	"github.com/alfianyulianto/pds-service/internal/model"
	"github.com/alfianyulianto/pds-service/internal/model/converter"
	"github.com/alfianyulianto/pds-service/internal/repository"
	"github.com/alfianyulianto/pds-service/internal/utils"
	"github.com/alfianyulianto/pds-service/pkg/email"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AccountUseCase interface {
	Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error)
	UpdatePassword(ctx context.Context, request *model.UpdatePasswordRequest) (*model.UserResponse, error)
}

type accountUseCase struct {
	*BaseUseCase
	UserRepository repository.UserRepository
	EmailService   *email.EmailService
}

func NewAccountUseCase(baseUseCase *BaseUseCase, userRepository repository.UserRepository, email *email.EmailService) AccountUseCase {
	return &accountUseCase{BaseUseCase: baseUseCase, UserRepository: userRepository, EmailService: email}
}

func (u *accountUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithField("action", "current").WithError(err).Warn("Failed to validate request body")
		return nil, err
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.ID); err != nil {
		u.Log.WithField("action", "current").WithError(err).Error("Failed to find user")
		return nil, fiber.NewError(fiber.StatusNotFound, "User data not found")
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithField("action", "current").WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (u *accountUseCase) UpdatePassword(ctx context.Context, request *model.UpdatePasswordRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithField("action", "update password").WithError(err).Warn("Failed to validate request body")
		return nil, err
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.ID); err != nil {
		u.Log.WithField("action", "update password").WithError(err).Error("Failed to find user")
		return nil, fiber.NewError(fiber.StatusNotFound, "User data not found")
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		u.Log.WithField("action", "update password").WithError(err).Error("Failed to hash password")
		return nil, fiber.ErrInternalServerError
	}
	user.Password = string(newPassword)

	if err = u.UserRepository.Update(tx, user); err != nil {
		u.Log.WithField("action", "update password").WithError(err).Error("Failed to update user password")
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		u.Log.WithField("action", "update password").WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	go func() {
		err = email.QuickPasswordChangedEmail(u.EmailService, user.Email, user.Name, utils.FormatTime(time.Now()))
		if err != nil {
			u.Log.WithField("action", "update password").WithError(err).Error("Failed to send password changed email")
		}
	}()

	return converter.UserToResponse(user), nil
}
