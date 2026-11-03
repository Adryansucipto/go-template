package partner

import (
	"context"
	"go-template/internal/controller/model"
)

type Repository interface {
	FetchUsersViaHTTP(ctx context.Context) ([]model.User, error)
}

