package impl

import (
	"context"
	"go-template/internal/repository/model"
	"go-template/util"
)

func (r *Repository) CreateRole(ctx context.Context, request model.Role) error {
	db := r.DbGorm.WithContext(ctx).Debug()
	if err := db.Create(&request).Error; err != nil {
		return util.ErrorGenerate(tag, err)
	}
	return nil
}
