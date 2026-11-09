package impl

import (
	"go-template/internal/repository"
	authService "go-template/internal/service/auth"
	"go-template/pkg/config"
)

var (
	tag = "[AuthService]"
)

type Service struct {
	Repository repository.Repository
	Config     config.Config
}

func New(repository repository.Repository, cfg config.Config) authService.Service {
	return &Service{
		Repository: repository,
		Config:     cfg,
	}
}
