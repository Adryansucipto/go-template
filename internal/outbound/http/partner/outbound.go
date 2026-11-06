package partner

import (
	"context"
	"go-template/internal/controller/model"
)

type Outbound interface {
	FetchUsersViaHTTP(ctx context.Context) ([]model.User, error)
}
