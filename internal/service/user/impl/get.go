package impl

import (
	"context"
	"fmt"
	"go-template/internal/controller/model"
	"go-template/internal/repository"
	"go-template/pkg/response"
)

func (s *Service) GetUsers(ctx context.Context) (res response.Response, err error) {
	var users []model.User

	nilRepo := repository.Repository{}
	if s.Repository == nilRepo {
		return response.ResponseGenerate(500, fmt.Errorf("repo is nil in service"), nil, nil, nil), fmt.Errorf("repo is nil in service")
	}

	// 1. Check Redis Cache first via Outbound layer (Cache-Aside Pattern)
	if s.Outbound.Cache.User != nil {
		if cachedUsers, hit, cacheErr := s.Outbound.Cache.User.GetUsersCache(ctx); cacheErr == nil && hit {
			return response.ResponseGenerate(200, nil, nil, cachedUsers, nil), nil
		}
	}

	// 2. Cache miss -> Query PostgreSQL via Repository layer
	records, err := s.Repository.UserRepository.GetUserRepository(ctx)
	if err != nil {
		return response.ResponseGenerate(500, fmt.Errorf("[GetUsers]-%s-%+v", tag, err.Error()), nil, nil, nil), fmt.Errorf("[GetUsers]-%s-%+v", tag, err.Error())
	}

	for _, v := range records {
		var user model.User
		user.ID = v.ID
		user.Name = v.Name
		user.Username = v.Username

		users = append(users, user)
	}

	// 3. Populate Redis Cache via Outbound layer
	if s.Outbound.Cache.User != nil {
		_ = s.Outbound.Cache.User.SetUsersCache(ctx, users)
	}

	return response.ResponseGenerate(200, nil, nil, users, nil), nil
}
