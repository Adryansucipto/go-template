package auth

import (
	"context"
	"go-template/internal/controller/model"
	repoModel "go-template/internal/repository/model"
	"go-template/pkg/response"
)

type Usecase interface {
	LoginHandler(ctx context.Context, request model.AuthRequest) (res response.Response)
	GetRecordByToken(ctx context.Context, token string) (repoModel.Session, response.Response)
	DeleteSession(ctx context.Context, record repoModel.Session) response.Response
	CleanupExpiredSessions(ctx context.Context) (int64, error)
}
