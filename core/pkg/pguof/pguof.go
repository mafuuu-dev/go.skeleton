package pguof

import (
	"backend/core/pkg/errorsx"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UnitOfWork struct {
	ctx context.Context
	db  *pgxpool.Pool
	log *zap.SugaredLogger
}

func New(ctx context.Context, db *pgxpool.Pool, log *zap.SugaredLogger) *UnitOfWork {
	return &UnitOfWork{
		ctx: ctx,
		db:  db,
		log: log,
	}
}

func (u *UnitOfWork) Do(fn func(tx pgx.Tx) error) error {
	tx, err := u.db.Begin(u.ctx)
	if err != nil {
		return errorsx.Wrap(err, "Failed to begin transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(u.ctx)

			panic(errorsx.Wrap(errorsx.Recover(r), "Panic recovered in UnitOfWork"))
		}

		if err := tx.Rollback(u.ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			u.log.Warn(errorsx.WrapJSON(err, "Rollback failed"))
		}
	}()

	if err := fn(tx); err != nil {
		return errorsx.Wrap(err, "Failed to execute transaction")
	}

	return tx.Commit(u.ctx)
}
