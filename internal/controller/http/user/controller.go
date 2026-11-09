package user

import (
	userService "go-template/internal/service/user"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

var (
	tag = "[UserController]"
)

type Controller struct {
	dig.In
	User   userService.Service
	Logger *zap.Logger
}
