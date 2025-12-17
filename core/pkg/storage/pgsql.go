package storage

import (
	"backend/core/config"
	"backend/core/pkg/errorsx"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errorsx.Error(err)
	}

	cfg.MaxConns = 100
	cfg.MinConns = 10
	cfg.MaxConnIdleTime = 30 * time.Second
	cfg.MaxConnLifetime = 10 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errorsx.Error(err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, errorsx.Error(err)
	}

	return pool, nil
}

func PostgresDsn(cfg *config.Config) string {
	return "postgres://" +
		cfg.PGUser + ":" +
		cfg.PGPassword + "@" +
		cfg.PGHost + ":" +
		cfg.PGPort + "/" +
		cfg.PGDatabase +
		"?sslmode=" + cfg.PGSSLMode
}
