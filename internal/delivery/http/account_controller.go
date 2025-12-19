package http

import (
	"github.com/alfianyulianto/pds-service/internal/delivery/http/middleware"
	"github.com/alfianyulianto/pds-service/internal/model"
	"github.com/alfianyulianto/pds-service/internal/usecase"
	"github.com/alfianyulianto/pds-service/pkg/auth"
	"github.com/alfianyulianto/pds-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AccountController interface {
	Current(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	UpdatePassword(ctx *fiber.Ctx) error
}
type accountController struct {
	UseCase    usecase.AccountUseCase
	Log        *logrus.Entry
	JwtService *auth.JWTService
}

func NewAccountController(useCase usecase.AccountUseCase, log *logrus.Entry, jwtService *auth.JWTService) AccountController {
	return &accountController{UseCase: useCase, Log: log, JwtService: jwtService}
}

func (c *accountController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.GetUserRequest{
		ID: auth.ID,
	}

	user, err := c.UseCase.Current(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[*model.UserResponse]{
		Success: true,
		Message: "User fetched successfully",
		Data:    user,
	})
}

func (c *accountController) Logout(ctx *fiber.Ctx) error {
	claims := middleware.GetUser(ctx)

	err := c.JwtService.RevokeToken(ctx.Context(), claims)
	if err != nil {
		c.Log.WithField("action", "logout").WithError(err).Error("Failed to revoke token")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to logout user")
	}

	return ctx.Status(200).JSON(response.Response[*model.UserResponse]{
		Success: true,
		Message: "User logged out successfully",
	})
}

func (c *accountController) UpdatePassword(ctx *fiber.Ctx) error {
	request := new(model.UpdatePasswordRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithField("action", "update password").WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	auth := middleware.GetUser(ctx)
	request.ID = auth.ID

	user, err := c.UseCase.UpdatePassword(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[*model.UserResponse]{
		Success: true,
		Message: "User password updated successfully",
		Data:    user,
	})
}
