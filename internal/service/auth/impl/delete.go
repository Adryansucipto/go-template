package impl

import (
	"context"
	repoModel "go-template/internal/repository/model"
	"go-template/pkg/response"
)

func (s *Service) DeleteSession(ctx context.Context, record repoModel.Session) response.Response {
	err := s.Repository.AuthRepository.DeleteSession(ctx, record.Username)
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

func (s *Service) CleanupExpiredSessions(ctx context.Context) (int64, error) {
	return s.Repository.AuthRepository.DeleteExpiredSessions(ctx)
}
