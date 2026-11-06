package impl

import (
	"context"
	repoModel "go-template/internal/repository/model"
	"go-template/internal/repository/transaction"
	"go-template/pkg/response"
)

func (r *Repository) CreateSession(ctx context.Context, request repoModel.Session) error {
	db := transaction.ExtractTx(ctx, r.DbGorm)
	if err := db.Create(&request).Error; err != nil {
		return response.ErrorGenerate(tag, err)
	}
	return nil
}
