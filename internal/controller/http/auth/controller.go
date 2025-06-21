package auth

import (
	authUsecase "go-template/internal/usecase/auth"

	"go.uber.org/dig"
)

var (
	tag = "[AuthController]"
)

type Controller struct {
	dig.In
	Auth authUsecase.Usecase
}
