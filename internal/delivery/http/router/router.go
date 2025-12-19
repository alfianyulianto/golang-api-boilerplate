package router

import (
	"github.com/alfianyulianto/pds-service/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"

	"github.com/alfianyulianto/pds-service/internal/delivery/http"
)

type RouterConfig struct {
	App               *fiber.App
	Middleware        *middleware.Middleware
	AuthController    http.AuthController
	AccountController http.AccountController
	UserController    http.UserController
}

func (c RouterConfig) Setup() {
	c.App.Static("/uploads", "./uploads")

	c.setupGuestRoute()
	c.setupAuthRoute()
}

func (c RouterConfig) setupGuestRoute() {
	auth := c.App.Group("/api/auth")
	auth.Post("/register", c.AuthController.Register)
	auth.Post("/login", c.AuthController.Login)
	auth.Post("/request-reset-password", c.AuthController.RequestResetPassword)
	auth.Post("/reset-password", c.AuthController.ResetPassword)
}

func (c RouterConfig) setupAuthRoute() {
	auth := c.App.Group("/api/auth", c.Middleware.AuthMiddleware)
	auth.Get("/_current", c.AccountController.Current)
	auth.Post("/refresh-token", c.AuthController.RefreshToken)
	auth.Post("/logout", c.AccountController.Logout)
	auth.Post("/update-password", c.AccountController.UpdatePassword)

	user := c.App.Group("/api/users", c.Middleware.AuthMiddleware)
	user.Get("/", c.UserController.List)
	user.Post("/", c.UserController.Create)
	user.Get("/:id", c.UserController.FindById)
	user.Put("/:id", c.UserController.Update)
	user.Delete("/:id", c.UserController.Delete)

}
