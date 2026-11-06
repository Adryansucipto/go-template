package user

import (
	"context"
	"go-template/internal/controller/model"
	"go-template/pkg/response"
)

type Service interface {
	GetUsers(ctx context.Context) (res response.Response, err error)
	CreateUser(ctx context.Context, req model.CreateUserRequest) (res response.Response, err error)
}
