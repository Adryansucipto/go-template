package impl

import (
	"go-template/internal/repository"
	userUsecase "go-template/internal/usecase/user"
)

var (
	tag = "[UserUsecase]"
)

type Usecase struct {
	Repository repository.Repository
}

func New(repository repository.Repository) userUsecase.Usecase {
	return &Usecase{
		Repository: repository,
	}
}
