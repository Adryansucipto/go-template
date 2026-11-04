package impl

import (
	"context"
	"go-template/internal/repository/db/tx"
	repoModel "go-template/internal/repository/model"
	"go-template/pkg/response"
)

func (r *Repository) DeleteSession(ctx context.Context, username string) error {
	db := tx.ExtractTx(ctx, r.DbGorm)

	if err := db.Where("username = ?", username).Delete(&repoModel.Session{}).Error; err != nil {
		return response.ErrorGenerate(tag, err)
	}

	return nil
}

func (r *Repository) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	db := tx.ExtractTx(ctx, r.DbGorm)

	res := db.Where("expired_session <= (now() at time zone 'Asia/Jakarta')").Delete(&repoModel.Session{})
	if res.Error != nil {
		return 0, response.ErrorGenerate(tag, res.Error)
	}

	return res.RowsAffected, nil
}
