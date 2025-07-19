package role

import (
	"context"
	"go-template/internal/controller/model"
	"go-template/util"
)

type Usecase interface {
	CreateRole(ctx context.Context, request model.CreateRoleRequest, username string) (util.Response, error)
}
