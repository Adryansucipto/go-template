package impl

import (
	"go-template/internal/repository"
	authUsecase "go-template/internal/usecase/auth"
	"go-template/pkg/config"

	"go.uber.org/dig"
)

var (
	tag = "[AuthUsecase]"
)

type Usecase struct {
	dig.In

	Repository repository.Repository
	Config     config.Config
}

func New(repository repository.Repository) authUsecase.Usecase {
	return &Usecase{
		Repository: repository,
	}
}
