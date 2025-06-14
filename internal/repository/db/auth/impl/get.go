package impl

import (
	"context"
	"fmt"
	"go-template/internal/repository/model"
)

func (r *Repository) GetUsernamePassword(ctx context.Context, username string) (record model.User, err error) {

	db := r.DbGorm.WithContext(ctx).Debug()

	if db == nil {
		fmt.Printf("LOL")
	}

	if err = db.Where("username = ?", username).First(&record).Error; err != nil {
		return model.User{}, err
	}

	return record, nil
}
