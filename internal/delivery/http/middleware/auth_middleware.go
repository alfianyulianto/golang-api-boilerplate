package middleware

import (
	"github.com/alfianyulianto/pds-service/internal/model"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) AuthMiddleware(ctx *fiber.Ctx) error {
	authHeader := ctx.Get("Authorization", "")
	if authHeader == "" {
		m.Log.WithField("action", "authentication middleware").Warn("Authorization header is missing")
		return fiber.ErrUnauthorized
	}

	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		m.Log.WithField("action", "authentication middleware").Warn("Invalid authorization header format")
		return fiber.ErrUnauthorized
	}

	request := &model.VerifyUserRequest{Token: token}
	m.Log.WithField("action", "authentication middleware").Debugf("Authorization : %s", request.Token)

	userClaim, err := m.Jwt.ParseAccessToken(ctx.Context(), request.Token)
	if err != nil {
		m.Log.WithField("action", "authentication middleware").WithError(err).Warn("Failed find user by token:", err)
		return fiber.ErrUnauthorized
	}

	m.Log.Debugf("Auth:", userClaim)
	ctx.Locals("auth", userClaim)
	return ctx.Next()
}

func GetUser(ctx *fiber.Ctx) *model.UserClaimToken {
	return ctx.Locals("auth").(*model.UserClaimToken)
}
