package impl

import (
	"context"
	"go-template/internal/repository/model"
	"go-template/internal/repository/transaction"
)

func (r *Repository) CreateUserRepository(ctx context.Context, user model.User) (model.User, error) {
	db := transaction.ExtractTx(ctx, r.DbGorm)

	if err := db.Create(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}
