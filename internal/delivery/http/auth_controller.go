package http

import (
	"context"
	"github.com/alfianyulianto/pds-service/internal/model"
	"github.com/alfianyulianto/pds-service/internal/usecase"
	"github.com/alfianyulianto/pds-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"strings"
)

type AuthController interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	RefreshToken(ctx *fiber.Ctx) error
	RequestResetPassword(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
}

type authController struct {
	UseCase usecase.AuthUseCase
	Log     *logrus.Entry
}

func NewAuthController(useCase usecase.AuthUseCase, log *logrus.Entry) AuthController {
	return &authController{UseCase: useCase, Log: log}
}

func (c *authController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithField("action", "register").WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := c.UseCase.Register(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[*model.UserResponse]{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

func (c *authController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithField("action", "login").WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	ua := strings.ToLower(ctx.Get("User-Agent"))
	device := "desktop"
	if strings.Contains(ua, "mobile") {
		device = "mobile"
	}

	deviceContext := context.WithValue(ctx.UserContext(), "DeviceTypeKey", device)
	ctx.SetUserContext(deviceContext)

	token, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[*model.AuthResponse]{
		Success: true,
		Message: "User logged in successfully",
		Data:    token,
	})
}

func (c *authController) RefreshToken(ctx *fiber.Ctx) error {
	request := new(model.VerifyUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithField("action", "refresh token").WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	token, err := c.UseCase.RefreshToken(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[*model.AuthResponse]{
		Success: true,
		Message: "Token refreshed successfully",
		Data:    token,
	})
}

func (c *authController) RequestResetPassword(ctx *fiber.Ctx) error {
	request := new(model.RequestPasswordResetRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithField("action", "request reset password").WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.UseCase.RequestResetPassword(ctx.Context(), request); err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[any]{
		Success: true,
		Message: "Password reset request successful. Please check your email.",
	})
}

func (c *authController) ResetPassword(ctx *fiber.Ctx) error {
	request := new(model.ResetPasswordRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithField("action", "reset password").WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := c.UseCase.ResetPassword(ctx.Context(), request); err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[any]{
		Success: true,
		Message: "Password has been reset successfully.",
	})
}
