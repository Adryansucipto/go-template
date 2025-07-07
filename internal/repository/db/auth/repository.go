package auth

import (
	"context"
	"go-template/internal/repository/model"
)

type Repository interface {
	GetUsernamePassword(ctx context.Context, username string) (record model.User, err error)
	GetActiveSession(ctx context.Context, username string) (record model.Session, err error)
	CreateSession(ctx context.Context, request model.Session) error
	GetActiveSessionByToken(ctx context.Context, token string) (record model.Session, err error)
	DeleteSession(ctx context.Context, username string) error
}
