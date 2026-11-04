package model

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Enabled    bool
	Client     *redis.Client
	DefaultTTL time.Duration
}

func (r *Redis) GetJSON(ctx context.Context, key string, dest any) (bool, error) {
	if !r.Enabled || r.Client == nil {
		return false, nil
	}

	val, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Redis) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	if !r.Enabled || r.Client == nil {
		return nil
	}

	if ttl <= 0 {
		ttl = r.DefaultTTL
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, key, bytes, ttl).Err()
}

func (r *Redis) Del(ctx context.Context, keys ...string) error {
	if !r.Enabled || r.Client == nil || len(keys) == 0 {
		return nil
	}
	return r.Client.Del(ctx, keys...).Err()
}

func (r *Redis) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}

