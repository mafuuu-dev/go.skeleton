package storage

import (
	"backend/core/config"
	"backend/core/pkg/errorsx"
	"backend/core/pkg/lifecycle"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Storage struct {
	PG    *pgxpool.Pool
	Redis *redis.Client
}

func New(ctx context.Context, cfg *config.Config, log *zap.SugaredLogger) *Storage {
	timeoutCtx, cancel := lifecycle.Timeout(ctx, 20*time.Second)
	defer cancel()

	database, err := Connect(timeoutCtx, cfg)
	if err != nil {
		log.Warn(errorsx.WrapJSON(err, "Failed to connect to storage"))
	}

	return database
}

func Connect(ctx context.Context, cfg *config.Config) (*Storage, error) {
	postgresClient, err := ConnectPostgres(ctx, PostgresDsn(cfg))
	if err != nil {
		return nil, errorsx.Wrap(err, "Failed to connect to postgres")
	}

	redisClient, err := ConnectRedis(ctx, RedisDsn(cfg))
	if err != nil {
		postgresClient.Close()
		return nil, errorsx.Wrap(err, "Failed to connect to redis")
	}

	return &Storage{
		PG:    postgresClient,
		Redis: redisClient,
	}, nil
}
