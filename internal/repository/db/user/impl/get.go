package impl

import (
	"context"
	"go-template/internal/repository/model"
)

func (r *Repository) GetUserRepository(ctx context.Context) (user []model.User, err error) {
	db := r.DbGorm.WithContext(ctx)

	if err = db.Find(&user).Error; err != nil {
		return []model.User{}, err
	}

	return user, nil
}
