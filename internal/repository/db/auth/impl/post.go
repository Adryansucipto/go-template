package impl

import (
	"context"
	"go-template/internal/repository/db/tx"
	repoModel "go-template/internal/repository/model"
	"go-template/pkg/response"
)

func (r *Repository) CreateSession(ctx context.Context, request repoModel.Session) error {
	db := tx.ExtractTx(ctx, r.DbGorm)
	if err := db.Create(&request).Error; err != nil {
		return response.ErrorGenerate(tag, err)
	}
	return nil
}
