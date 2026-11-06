package injection

import (
	pkgConfig "go-template/pkg/config"
	"go-template/pkg/resource/model"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg pkgConfig.Config) (model.Redis, error) {
	defaultTTL := time.Duration(cfg.Redis.DefaultTTLMin) * time.Minute
	if defaultTTL <= 0 {
		defaultTTL = 10 * time.Minute
	}

	if !cfg.Redis.Enabled {
		return model.Redis{
			Enabled:    false,
			DefaultTTL: defaultTTL,
		}, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return model.Redis{
		Enabled:    true,
		Client:     client,
		DefaultTTL: defaultTTL,
	}, nil
}
