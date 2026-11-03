package impl

import (
	"context"
	repoModel "go-template/internal/repository/model"
	"go-template/pkg/response"
)

func (u *Usecase) DeleteSession(ctx context.Context, record repoModel.Session) response.Response {
	err := u.Repository.DBRepository.AuthRepository.DeleteSession(ctx, record.Username)
	if err != nil {
		return response.ResponseGenerate(
			500,
			err,
			nil,
			nil,
			nil,
		)
	}
	return response.ResponseGenerate(
		200,
		nil,
		nil,
		nil,
		nil,
	)
}

func (u *Usecase) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	return u.Repository.DBRepository.AuthRepository.DeleteExpiredSessions(ctx)
}
