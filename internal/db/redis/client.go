package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"gitlab.snapp.ir/dispatching/soteria/configs"
)

func NewRedisClient(cfg *configs.RedisConfig) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:               cfg.Address,
		Password:           cfg.Password,
		PoolTimeout:        cfg.PoolTimeout,
		ReadTimeout:        cfg.ReadTimeout,
		MaxRetries:         cfg.MaxRetries,
		MaxRetryBackoff:    cfg.MaxRetryBackoff,
		MinRetryBackoff:    cfg.MinRetryBackoff,
		IdleTimeout:        cfg.IdleTimeout,
		IdleCheckFrequency: cfg.IdleCheckFrequency,
		PoolSize:           cfg.PoolSize,
	}

	client := redis.NewClient(opts)
	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("could not ping redis: %w", err)
	}

	return client, nil
}
