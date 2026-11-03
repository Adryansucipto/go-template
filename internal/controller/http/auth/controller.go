package auth

import (
	authUsecase "go-template/internal/usecase/auth"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

var (
	tag = "[AuthController]"
)

type Controller struct {
	dig.In
	Auth   authUsecase.Usecase
	Logger *zap.Logger
}
