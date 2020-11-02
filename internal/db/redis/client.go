package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
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
		MinIdleConns:       cfg.MinIdleConnections,
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("could not ping redis: %w", err)
	}

	return client, nil
}
