package impl

import (
	"context"
	"errors"
	"go-template/internal/repository/model"
	"go-template/pkg/response"
	"reflect"
)

func (s *Service) CheckTokenExpired(ctx context.Context, token, username string) (int, error) {
	record, err := s.Repository.AuthRepository.GetActiveSession(ctx, username)
	if err != nil {
		return 500, response.ErrorGenerate(
			tag,
			err,
		)
	}

	if !reflect.DeepEqual(record, model.Session{}) {
		return 400, response.ErrorGenerate(
			tag,
			errors.New("session still active, token not expired"),
		)
	}

	return 200, nil
}
