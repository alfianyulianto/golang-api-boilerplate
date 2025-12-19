package usecase

import (
	"context"
	"fmt"
	"github.com/alfianyulianto/pds-service/internal/entity"
	"github.com/alfianyulianto/pds-service/internal/model"
	"github.com/alfianyulianto/pds-service/internal/model/converter"
	"github.com/alfianyulianto/pds-service/internal/repository"
	"github.com/alfianyulianto/pds-service/internal/utils"
	"github.com/alfianyulianto/pds-service/pkg/auth"
	"github.com/alfianyulianto/pds-service/pkg/email"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthUseCase interface {
	Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.AuthResponse, error)
	RefreshToken(ctx context.Context, request *model.VerifyUserRequest) (*model.AuthResponse, error)
	RequestResetPassword(ctx context.Context, request *model.RequestPasswordResetRequest) error
	ResetPassword(ctx context.Context, request *model.ResetPasswordRequest) error
}

type authUseCase struct {
	*BaseUseCase
	UserRepository repository.UserRepository
	JwtService     *auth.JWTService
	EmailService   *email.EmailService
	Redis          *redis.Client
}

func NewAuthUseCase(baseUseCase *BaseUseCase, userRepository repository.UserRepository, jwtService *auth.JWTService, emailService *email.EmailService, redis *redis.Client) AuthUseCase {
	return &authUseCase{BaseUseCase: baseUseCase, UserRepository: userRepository, JwtService: jwtService, EmailService: emailService, Redis: redis}
}

func (u *authUseCase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithField("action", "login").WithError(err).Warn("Failed to validate request body")
		return nil, err
	}

	user := converter.RegisterRequestToUser(request)

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		u.Log.WithField("action", "login").WithError(err).Error("Failed to hash password")
		return nil, fiber.ErrInternalServerError
	}

	user.Password = string(password)

	if err = u.UserRepository.Create(tx, user); err != nil {
		u.Log.WithField("action", "login").WithError(err).Error("Failed to register user")
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		u.Log.WithField("action", "login").WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	go func() {
		err = email.QuickSendWelcome(u.EmailService, user.Email, user.Name)
	}()

	return converter.UserToResponse(user), nil
}

func (u *authUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.AuthResponse, error) {
	device := ctx.Value("DeviceTypeKey").(string)

	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithField("action", "login").WithError(err).Warn("Failed to validate request body")
		return nil, err
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		u.Log.WithField("action", "login").WithError(err).Error("Failed to find user by email")
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		u.Log.WithField("action", "login").WithError(err).Error("Invalid password")
		return nil, fiber.ErrUnauthorized
	}

	claims := model.UserClaimToken{
		ID:   user.ID,
		Role: user.Role,
	}
	token, err := u.JwtService.CreateToken(ctx, &claims)
	if err != nil {
		u.Log.WithField("action", "login").WithError(err).Error("Failed to create JWT token")
		return nil, fiber.ErrInternalServerError
	}

	lastLogIntAt := time.Now()
	user.LastLoginAt = &lastLogIntAt
	if err = u.UserRepository.Update(tx, user); err != nil {
		u.Log.WithField("action", "login").WithError(err).Error("Failed to update user last login")
		return nil, fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		u.Log.WithField("action", "login").WithError(err).Error("Failed to commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	go func() {
		err = email.QuickSendLoginNotification(u.EmailService, user.Email, user.Name, utils.FormatTime(lastLogIntAt), device)
		if err != nil {
			u.Log.WithField("action", "login").WithError(err).Error("Failed to send login notification email")
		}
	}()

	return token, nil
}

func (u *authUseCase) RefreshToken(ctx context.Context, request *model.VerifyUserRequest) (*model.AuthResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithField("action", "refresh token").WithError(err).Warn("Failed to validate request body")
		return nil, err
	}

	claims, err := u.JwtService.ParseRefreshToken(ctx, request.Token)
	if err != nil {
		u.Log.WithField("action", "refresh token").WithError(err).Error("Failed to parse access token")
		return nil, fiber.ErrUnauthorized
	}

	token, err := u.JwtService.CreateToken(ctx, claims)
	if err != nil {
		u.Log.WithField("action", "refresh token").WithError(err).Error("Failed to create new access token, from refresh token")
		return nil, fiber.ErrInternalServerError
	}

	return token, nil
}

func (u *authUseCase) RequestResetPassword(ctx context.Context, request *model.RequestPasswordResetRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithField("action", "request reset password").WithError(err).Warn("Failed to validate request body")
		return err
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		u.Log.WithField("action", "request reset password").WithError(err).Error("Failed to find user by email")
		return fiber.NewError(fiber.StatusNotFound, "User data not found")
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithField("action", "request reset password").WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	token := uuid.NewString()
	err := u.Redis.SetEx(ctx, "reset_password:"+token, user.Email, 1*time.Hour).Err()
	if err != nil {
		u.Log.WithField("action", "request reset password").WithError(err).Error("Failed to set reset password token in redis")
		return fiber.ErrInternalServerError
	}

	go func() {
		resetURL := fmt.Sprintf("https://alfian.my.id/accounts/password/reset/confirm?token=%s", token)
		err = email.QuickSendResetPassword(u.EmailService, user.Email, user.Name, resetURL)
		if err != nil {
			u.Log.WithField("action", "request reset password").WithError(err).Error("Failed to send reset password email")
		}
	}()

	return nil
}

func (u *authUseCase) ResetPassword(ctx context.Context, request *model.ResetPasswordRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithField("action", "reset password").WithError(err).Warn("Failed to validate request body")
		return err
	}

	userEmail, err := u.Redis.Get(ctx, "reset_password:"+request.Token).Result()
	if err != nil {
		u.Log.WithField("action", "reset password").WithError(err).Error("Failed to get reset password token from redis")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid or expired reset password token")
	}

	user := new(entity.User)
	if err = u.UserRepository.FindByEmail(tx, user, userEmail); err != nil {
		u.Log.WithField("action", "reset password").WithError(err).Error("Failed to find user by email")
		return fiber.NewError(fiber.StatusNotFound, "User data not found")
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		u.Log.WithField("action", "reset password").WithError(err).Error("Failed to hash new password")
		return fiber.ErrInternalServerError
	}
	user.Password = string(newPassword)

	if err = u.UserRepository.Update(tx, user); err != nil {
		u.Log.WithField("action", "reset password").WithError(err).Error("Failed to update user password")
		return fiber.ErrInternalServerError
	}

	if err = tx.Commit().Error; err != nil {
		u.Log.WithField("action", "reset password").WithError(err).Error("Failed to commit transaction")
		return fiber.ErrInternalServerError
	}

	err = u.Redis.Del(ctx, "reset_password:"+request.Token).Err()
	if err != nil {
		u.Log.WithField("action", "reset password").WithError(err).Error("Failed to delete reset password token from redis")
		return fiber.ErrInternalServerError
	}

	go func() {
		err = email.QuickPasswordChangedEmail(u.EmailService, user.Email, user.Name, utils.FormatTime(time.Now()))
		if err != nil {
			u.Log.WithField("action", "reset password").WithError(err).Error("Failed to send password changed notification email")
		}
	}()

	return nil
}
