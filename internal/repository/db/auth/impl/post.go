package impl

import (
	"context"
	repoModel "go-template/internal/repository/model"
	"go-template/pkg/response"
)

func (r *Repository) CreateSession(ctx context.Context, request repoModel.Session) error {
	db := r.DbGorm.WithContext(ctx)
	if err := db.Create(&request).Error; err != nil {
		return response.ErrorGenerate(tag, err)
	}
	return nil
}
