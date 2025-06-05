package user

import (
	"context"
	"go-template/internal/controller/model"
)

type Usecase interface {
	GetUsers(ctx context.Context) ([]model.User, error)
}
