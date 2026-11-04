package impl

import (
	"context"
	"go-template/internal/controller/model"
	userCache "go-template/internal/repository/cache/user"
	resModel "go-template/pkg/resource/model"
	"time"
)

const usersCacheKey = "users:list:all"

type Repository struct {
	Redis resModel.Redis
}

func New(redis resModel.Redis) userCache.Repository {
	return &Repository{
		Redis: redis,
	}
}

func (r *Repository) GetUsersCache(ctx context.Context) ([]model.User, bool, error) {
	var users []model.User
	hit, err := r.Redis.GetJSON(ctx, usersCacheKey, &users)
	if err != nil || !hit {
		return nil, false, err
	}
	return users, true, nil
}

func (r *Repository) SetUsersCache(ctx context.Context, users []model.User) error {
	return r.Redis.SetJSON(ctx, usersCacheKey, users, 5*time.Minute)
}

func (r *Repository) InvalidateUsersCache(ctx context.Context) error {
	return r.Redis.Del(ctx, usersCacheKey)
}

