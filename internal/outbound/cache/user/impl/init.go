package impl

import (
	"context"
	"go-template/internal/controller/model"
	userCache "go-template/internal/outbound/cache/user"
	resModel "go-template/pkg/resource/model"
	"time"
)

const usersCacheKey = "users:list:all"

type Outbound struct {
	Redis resModel.Redis
}

func New(redis resModel.Redis) userCache.Outbound {
	return &Outbound{
		Redis: redis,
	}
}

func (o *Outbound) GetUsersCache(ctx context.Context) ([]model.User, bool, error) {
	var users []model.User
	hit, err := o.Redis.GetJSON(ctx, usersCacheKey, &users)
	if err != nil || !hit {
		return nil, false, err
	}
	return users, true, nil
}

func (o *Outbound) SetUsersCache(ctx context.Context, users []model.User) error {
	return o.Redis.SetJSON(ctx, usersCacheKey, users, 5*time.Minute)
}

func (o *Outbound) InvalidateUsersCache(ctx context.Context) error {
	return o.Redis.Del(ctx, usersCacheKey)
}
