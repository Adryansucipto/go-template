package impl

import (
	"go-template/internal/repository"
	roleUsecase "go-template/internal/usecase/role"
	"go-template/pkg/config"

	"go.uber.org/dig"
)

var (
	tag = "[RoleUsecase]"
)

type Usecase struct {
	dig.In

	Repository repository.Repository
	Config     config.Config
}

func New(repository repository.Repository) roleUsecase.Usecase {
	return &Usecase{
		Repository: repository,
	}
}
