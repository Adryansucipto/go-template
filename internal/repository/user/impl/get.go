package impl

import (
	"context"
	"go-template/internal/repository/model"
	"go-template/internal/repository/transaction"
)

func (r *Repository) GetUserRepository(ctx context.Context) (user []model.User, err error) {
	db := transaction.ExtractTx(ctx, r.DbGorm)

	if err = db.Find(&user).Error; err != nil {
		return []model.User{}, err
	}

	return user, nil
}
