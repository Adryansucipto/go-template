package user

import (
	userUsecase "go-template/internal/usecase/user"

	"go.uber.org/dig"
)

var (
	tag = "[UserController]"
)

type Controller struct {
	dig.In
	User userUsecase.Usecase
}
