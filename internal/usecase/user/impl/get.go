package impl

import (
	"context"
	"fmt"
	"go-template/internal/controller/model"
	"go-template/internal/repository"
	"go-template/pkg/response"
)

func (u *Usecase) GetUsers(ctx context.Context) (res response.Response, err error) {
	var users []model.User

	nilfunc := repository.Repository{}
	if u.Repository == nilfunc {
		return response.ResponseGenerate(500, fmt.Errorf("repo is nil in usecase"), nil, nil, nil), fmt.Errorf("repo is nil in usecase")
	}

	// 1. Check Redis Cache first (Cache-Aside Pattern)
	if u.Repository.CacheRepository.UserCache != nil {
		if cachedUsers, hit, cacheErr := u.Repository.CacheRepository.UserCache.GetUsersCache(ctx); cacheErr == nil && hit {
			return response.ResponseGenerate(200, nil, nil, cachedUsers, nil), nil
		}
	}

	// 2. Cache miss -> Query PostgreSQL
	records, err := u.Repository.DBRepository.UserRepository.GetUserRepository(ctx)
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

	// 3. Populate Redis Cache for subsequent requests
	if u.Repository.CacheRepository.UserCache != nil {
		_ = u.Repository.CacheRepository.UserCache.SetUsersCache(ctx, users)
	}

	return response.ResponseGenerate(200, nil, nil, users, nil), nil
}
