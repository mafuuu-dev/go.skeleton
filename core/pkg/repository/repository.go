package repository

import (
	"backend/core/pkg/query"
	"backend/core/types"
	"context"
)

type Repository struct {
	ctx context.Context
	db  types.PgExecutor
}

func NewRepository(ctx context.Context, db types.PgExecutor) *Repository {
	return &Repository{
		ctx: ctx,
		db:  db,
	}
}

func (repository *Repository) Context() context.Context {
	return repository.ctx
}

func (repository *Repository) DB() types.PgExecutor {
	return repository.db
}

func (repository *Repository) Query() *query.Factory {
	return query.New(repository.ctx, repository.db)
}
