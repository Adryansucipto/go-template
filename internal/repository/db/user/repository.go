package user

import (
	"context"
	"go-template/internal/repository/model"
)

type Repository interface {
	GetUserRepository(ctx context.Context) (user []model.User, err error)
}
