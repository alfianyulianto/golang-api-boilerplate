package middleware

import (
	"github.com/alfianyulianto/pds-service/pkg/auth"
	"github.com/sirupsen/logrus"
)

type Middleware struct {
	Log *logrus.Entry
	Jwt *auth.JWTService
}

func NewMiddleware(log *logrus.Entry, jwt *auth.JWTService) *Middleware {
	return &Middleware{Log: log, Jwt: jwt}
}
