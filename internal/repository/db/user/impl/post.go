package impl

import (
	"context"
	"go-template/internal/repository/model"
)

func (r *Repository) CreateUserRepository(ctx context.Context, user model.User) (model.User, error) {
	db := r.DbGorm.WithContext(ctx)

	if err := db.Create(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

