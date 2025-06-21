package impl

import (
	"context"
	repoModel "go-template/internal/repository/model"
	"go-template/util"
)

func (r *Repository) CreateSession(ctx context.Context, request repoModel.Session) error {
	db := r.DbGorm.WithContext(ctx).Debug()

	if err := db.Where("username = ?", request.Username).Delete(&repoModel.Session{}).Error; err != nil {
		return util.ErrorGenerate(tag, err)
	}

	if err := db.Create(request).Error; err != nil {
		return util.ErrorGenerate(tag, err)
	}
	return nil
}
