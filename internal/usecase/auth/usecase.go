package auth

import (
	"context"
	"go-template/internal/controller/model"
	"go-template/util"
)

type Usecase interface {
	LoginHandler(ctx context.Context, request model.AuthRequest) (response util.Response)
}
