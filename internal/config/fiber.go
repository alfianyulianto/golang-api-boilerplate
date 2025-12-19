package config

import (
	"github.com/alfianyulianto/pds-service/pkg/response"
	"github.com/alfianyulianto/pds-service/pkg/validators"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ErrorHandler: newErrorHandler,
		BodyLimit:    10 * 1024 * 1024,
	})

	return app
}

func newErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		return ctx.Status(400).JSON(response.Response[any]{
			Success: false,
			Message: "Validation Error",
			Error:   validators.ParseErrors(validationErrors),
		})
	}

	return ctx.Status(code).JSON(response.Response[any]{
		Success: false,
		Message: fiber.NewError(code).Message,
		Error:   err.Error(),
	})
}
