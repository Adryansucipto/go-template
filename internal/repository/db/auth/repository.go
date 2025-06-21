package auth

import (
	"context"
	"go-template/internal/repository/model"
)

type Repository interface {
	GetUsernamePassword(ctx context.Context, username string) (record model.User, err error)
	GetActiveSession(ctx context.Context, username string) bool
	CreateSession(ctx context.Context, request model.Session) error
}
