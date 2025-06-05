package user

import (
	userUsecase "go-template/internal/usecase/user"

	"go.uber.org/dig"
)

type Controller struct {
	dig.In
	User userUsecase.Usecase
}
