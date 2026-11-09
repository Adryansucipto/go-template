package impl

import (
	"go-template/internal/outbound"
	"go-template/internal/repository"
	userService "go-template/internal/service/user"
)

var (
	tag = "[UserService]"
)

type Service struct {
	Repository repository.Repository
	Outbound   outbound.Outbound
}

func New(repository repository.Repository, out outbound.Outbound) userService.Service {
	return &Service{
		Repository: repository,
		Outbound:   out,
	}
}
