package role

import (
	"context"
	"go-template/internal/repository/model"
)

type Repository interface {
	CreateRole(ctx context.Context, request model.Role) error
}
