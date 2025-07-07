package impl

import (
	"context"
	repoModel "go-template/internal/repository/model"
	"go-template/util"
)

func (r *Repository) DeleteSession(ctx context.Context, username string) error {
	db := r.DbGorm.WithContext(ctx).Debug()

	if err := db.Where("username = ?", username).Delete(&repoModel.Session{}).Error; err != nil {
		return util.ErrorGenerate(tag, err)
	}

	return nil
}
