package session

import (
	"context"
	authUsecase "go-template/internal/usecase/auth"

	"go.uber.org/dig"
	"go.uber.org/zap"
)

type CleanupJob struct {
	dig.In

	Auth   authUsecase.Usecase
	Logger *zap.Logger
}

func (j *CleanupJob) Run(ctx context.Context) {
	deletedCount, err := j.Auth.CleanupExpiredSessions(ctx)
	if err != nil {
		j.Logger.Error("[SessionCleanupJob] failed cleaning up expired sessions", zap.Error(err))
		return
	}
	if deletedCount > 0 {
		j.Logger.Info("[SessionCleanupJob] cleaned up expired sessions", zap.Int64("deleted_rows", deletedCount))
	}
}

