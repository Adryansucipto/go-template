package impl

import (
	"context"
	"fmt"
	"go-template/internal/repository/model"
)

func (r *Repository) GetUserRepository(ctx context.Context) (user []model.User, err error) {

	db := r.DbGorm.WithContext(ctx).Debug()

	if db == nil {
		fmt.Printf("LOL")
	}

	if err = db.Find(&user).Error; err != nil {
		return []model.User{}, err
	}

	return user, nil
}

func (r *Repository) GetUserByIDRepository(ctx context.Context, username string) (user model.User, err error) {

	db := r.DbGorm.WithContext(ctx).Debug()

	if db == nil {
		fmt.Printf("LOL")
	}

	if err = db.Select("id,name,username").Where("username = ?", username).Find(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}
