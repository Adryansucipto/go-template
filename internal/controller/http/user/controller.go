package user

import (
	userUsecase "go-template/internal/usecase/user"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

var (
	tag = "[UserController]"
)

type Controller struct {
	dig.In
	User   userUsecase.Usecase
	Logger *zap.Logger
}
