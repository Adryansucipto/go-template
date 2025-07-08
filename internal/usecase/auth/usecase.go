package auth

import (
	"context"
	"go-template/internal/controller/model"
	"go-template/util"
)

type Usecase interface {
	LoginHandler(ctx context.Context, request model.AuthRequest) (response util.Response)
	ValidateToken(ctx context.Context, token string) (string, util.Response)
	ValidateRefreshToken(ctx context.Context, token string) (string, util.Response)
	DeleteSession(ctx context.Context, username string) util.Response
	AuthorizeUser(ctx context.Context, token string) (model.User, util.Response)
	RefreshNewToken(ctx context.Context, username string) (response util.Response)
}
