package impl

import (
	"context"
	repoModel "go-template/internal/repository/model"
	"go-template/pkg/response"
)

func (r *Repository) DeleteSession(ctx context.Context, username string) error {
	db := r.DbGorm.WithContext(ctx)

	if err := db.Where("username = ?", username).Delete(&repoModel.Session{}).Error; err != nil {
		return response.ErrorGenerate(tag, err)
	}

	return nil
}

func (r *Repository) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	db := r.DbGorm.WithContext(ctx)

	tx := db.Where("expired_session <= (now() at time zone 'Asia/Jakarta')").Delete(&repoModel.Session{})
	if tx.Error != nil {
		return 0, response.ErrorGenerate(tag, tx.Error)
	}

	return tx.RowsAffected, nil
}
