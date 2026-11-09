package user

import (
	"context"
	"go-template/internal/controller/model"
)

type Outbound interface {
	GetUsersCache(ctx context.Context) ([]model.User, bool, error)
	SetUsersCache(ctx context.Context, users []model.User) error
	InvalidateUsersCache(ctx context.Context) error
}
