package storage

import (
	"backend/core/config"
	"backend/core/pkg/errorsx"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context, addr string) (*redis.Client, error) {
	options, err := redis.ParseURL(addr)
	if err != nil {
		return nil, errorsx.Error(err)
	}

	options.PoolSize = 100
	options.MinIdleConns = 20
	options.PoolTimeout = 1 * time.Second
	options.ReadTimeout = 1 * time.Second
	options.WriteTimeout = 1 * time.Second

	client := redis.NewClient(options)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errorsx.Error(err)
	}

	return client, nil
}

func RedisDsn(cfg *config.Config) string {
	return "redis://" +
		cfg.RedisHost + ":" +
		cfg.RedisPort + "/" +
		cfg.RedisDatabase
}
