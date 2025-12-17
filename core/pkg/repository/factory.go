package repository

import (
	"backend/core/types"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Factory struct {
	ctx context.Context
	db  types.PgExecutor
}

func New(ctx context.Context, db *pgxpool.Pool) *Factory {
	return &Factory{
		ctx: ctx,
		db:  db,
	}
}

func (f *Factory) WithTx(tx pgx.Tx) *Factory {
	return &Factory{
		ctx: f.ctx,
		db:  tx,
	}
}

func (f *Factory) Instance() *Repository {
	return NewRepository(f.ctx, f.db)
}
