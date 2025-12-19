package http

import (
	"github.com/alfianyulianto/pds-service/internal/model"
	"github.com/alfianyulianto/pds-service/internal/usecase"
	"github.com/alfianyulianto/pds-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
	FindById(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type userController struct {
	UseCase usecase.UserUseCase
	Log     *logrus.Entry
}

func NewUserController(useCase usecase.UserUseCase, log *logrus.Entry) UserController {
	return &userController{UseCase: useCase, Log: log}
}

func (c *userController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithField("action", "create user").WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var err error
	request.Avatar, err = ctx.FormFile("avatar")
	if err != nil {
		c.Log.WithField("action", "create user").WithError(err).Warn("Failed to parse avatar, because avatar is optional")
	}

	user, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[*model.UserResponse]{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

func (c *userController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithField("action", "update user").WithError(err).Error("Failed to parse request body")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var err error
	request.Avatar, err = ctx.FormFile("avatar")
	if err != nil {
		c.Log.WithField("action", "update user").WithError(err).Warn("Failed to parse avatar, because avatar is optional")
	}

	id := ctx.Params("id")
	request.ID, err = uuid.Parse(id)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to parse id")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err := c.UseCase.Update(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[*model.UserResponse]{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

func (c *userController) List(ctx *fiber.Ctx) error {
	request := new(model.SearchUserRequest)

	request.Search = ctx.Query("search")
	request.IsActive = ctx.Query("is_active")
	request.Role = ctx.Query("role")
	request.OrderBy = ctx.Query("order_by", "created_at")
	request.OrderDir = ctx.Query("order_dir", "desc")
	request.Page = ctx.QueryInt("page", 1)
	request.PageSize = ctx.QueryInt("page_size", 10)

	users, pagination, err := c.UseCase.List(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[*[]model.UserResponse]{
		Success:    true,
		Message:    "Users retrieved successfully",
		Data:       users,
		Pagination: pagination,
	})
}

func (c *userController) FindById(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.UseCase.FindById(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[*model.UserResponse]{
		Success: true,
		Message: "User data retrieved successfully",
		Data:    user,
	})
}

func (c *userController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.UseCase.Delete(ctx.Context(), id)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(response.Response[*model.UserResponse]{
		Success: true,
		Message: "User deleted successfully",
		Data:    nil,
	})
}
