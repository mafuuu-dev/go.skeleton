package query

import (
	"backend/core/types"
	"context"
)

type Factory struct {
	ctx context.Context
	db  types.PgExecutor
}

func New(ctx context.Context, db types.PgExecutor) *Factory {
	return &Factory{
		ctx: ctx,
		db:  db,
	}
}

func (f *Factory) Instance() *Query {
	return NewQuery(f.ctx, f.db)
}
