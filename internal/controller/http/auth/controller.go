package auth

import (
	authUsecase "go-template/internal/usecase/auth"

	"go.uber.org/dig"
)

type Controller struct {
	dig.In
	Auth authUsecase.Usecase
}
