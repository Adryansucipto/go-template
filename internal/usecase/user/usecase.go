package user

import (
	"context"
	"go-template/util"
)

type Usecase interface {
	GetUsers(ctx context.Context) (res util.Response, err error)
}
