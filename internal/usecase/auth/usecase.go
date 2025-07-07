package auth

import (
	"context"
	"go-template/internal/controller/model"
	repoModel "go-template/internal/repository/model"
	"go-template/util"
)

type Usecase interface {
	LoginHandler(ctx context.Context, request model.AuthRequest) (response util.Response)
	GetRecordByToken(ctx context.Context, token string) (repoModel.Session, util.Response)
	DeleteSession(ctx context.Context, record repoModel.Session) util.Response
}
