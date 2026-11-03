package impl

import (
	"go-template/internal/repository"
	authUsecase "go-template/internal/usecase/auth"
	"go-template/pkg/config"
)

var (
	tag = "[AuthUsecase]"
)

type Usecase struct {
	Repository repository.Repository
	Config     config.Config
}

func New(repository repository.Repository, cfg config.Config) authUsecase.Usecase {
	return &Usecase{
		Repository: repository,
		Config:     cfg,
	}
}
