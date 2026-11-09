package auth

import (
	authService "go-template/internal/service/auth"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

var (
	tag = "[AuthController]"
)

type Controller struct {
	dig.In
	Auth   authService.Service
	Logger *zap.Logger
}
