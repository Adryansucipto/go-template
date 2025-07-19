package role

import (
	authUsecase "go-template/internal/usecase/auth"
	roleUsecase "go-template/internal/usecase/role"

	"go.uber.org/dig"
)

var (
	tag = "[RoleController]"
)

type Controller struct {
	dig.In
	// Usecase
	Auth authUsecase.Usecase
	Role roleUsecase.Usecase
}
