package impl

import userUsecase "go-template/internal/usecase/user"

var (
	tag = "[UserUsecase]"
)

type Usecase struct {
}

func New() userUsecase.Usecase {
	return &Usecase{}
}
