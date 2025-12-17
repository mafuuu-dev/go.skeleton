package query

import (
	"backend/core/types"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SqlProvider interface {
	Sql() string
}

type Executor interface {
	Execute() error
}

type ExecutorResult[T any] interface {
	Execute() (T, error)
}

type Query struct {
	ctx context.Context
	db  types.PgExecutor
}

func NewQuery(ctx context.Context, db types.PgExecutor) *Query {
	return &Query{
		ctx: ctx,
		db:  db,
	}
}

func (query *Query) Exec(provider SqlProvider, args ...any) (pgconn.CommandTag, error) {
	return query.db.Exec(query.ctx, provider.Sql(), args...)
}

func (query *Query) QueryAll(provider SqlProvider, args ...any) (pgx.Rows, error) {
	return query.db.Query(query.ctx, provider.Sql(), args...)
}

func (query *Query) QueryRow(provider SqlProvider, args ...any) pgx.Row {
	return query.db.QueryRow(query.ctx, provider.Sql(), args...)
}

func (query *Query) Sql() string {
	panic("Sql() must be implemented")
}
