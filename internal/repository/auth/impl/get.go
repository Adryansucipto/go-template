package impl

import (
	"context"
	"errors"
	"go-template/internal/repository/model"
	"go-template/internal/repository/transaction"
	"go-template/pkg/response"

	"gorm.io/gorm"
)

func (r *Repository) GetUsernamePassword(ctx context.Context, username string) (record model.User, err error) {
	db := transaction.ExtractTx(ctx, r.DbGorm)

	if err = db.Where("username = ?", username).First(&record).Error; err != nil {
		return model.User{}, err
	}

	return record, nil
}

func (r *Repository) GetActiveSession(ctx context.Context, username string) (record model.Session, err error) {
	db := transaction.ExtractTx(ctx, r.DbGorm)

	if err := db.Select("expired_session,username").Where("username = ? AND expired_session > (now() at time zone 'Asia/Jakarta') ", username).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Session{}, nil
		}
		return model.Session{}, response.ErrorGenerate(tag, err)
	}
	return record, nil
}

func (r *Repository) GetActiveSessionByToken(ctx context.Context, token string) (record model.Session, err error) {
	db := transaction.ExtractTx(ctx, r.DbGorm)

	if err := db.Where("access_token = ? AND expired_session > (now() at time zone 'Asia/Jakarta')", token).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Session{}, nil
		}
		return model.Session{}, response.ErrorGenerate(tag, err)
	}
	return record, nil
}
