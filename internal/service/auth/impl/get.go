package impl

import (
	"context"
	"errors"
	"go-template/internal/repository/model"
	"go-template/pkg/response"
	"reflect"
	"strings"
)

func (s *Service) GetRecordByToken(ctx context.Context, token string) (model.Session, response.Response) {
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return model.Session{},
			response.ResponseGenerate(
				401,
				errors.New("missing or invalid bearer token"),
				nil,
				nil,
				nil,
			)
	}

	rawToken := strings.TrimPrefix(token, "Bearer ")
	record, err := s.Repository.AuthRepository.GetActiveSessionByToken(ctx, rawToken)
	if err != nil {
		return model.Session{},
			response.ResponseGenerate(
				500,
				err,
				nil,
				nil,
				nil,
			)
	}

	if reflect.DeepEqual(record, model.Session{}) {
		return model.Session{},
			response.ResponseGenerate(
				401,
				errors.New("session expired or token not found"),
				nil,
				nil,
				nil,
			)
	}

	return record, response.ResponseGenerate(200, nil, nil, nil, nil)
}
