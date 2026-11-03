package partner

import (
	"context"
	userv1 "go-template/proto/user/v1"
)

type Repository interface {
	FetchUsersViaGRPC(ctx context.Context) ([]*userv1.UserItem, error)
}

