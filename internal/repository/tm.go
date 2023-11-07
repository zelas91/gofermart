package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/zelas91/gofermart/internal/logger"
	"github.com/zelas91/gofermart/internal/types"
)

type transactionManager interface {
	do(ctx context.Context, fn func(ctx context.Context) error) error
	getConn(ctx context.Context) conn
}

type tm struct {
	db *sqlx.DB
}

func newTm(db *sqlx.DB) transactionManager {
	return &tm{
		db: db,
	}
}
func (t *tm) do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(sql.ErrTxDone, err) {
			logger.GetLogger(ctx).Errorf("Rollback err: %v", err)
			return
		}
	}()

	if err = fn(context.WithValue(ctx, types.TxKey, tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func (t *tm) getConn(ctx context.Context) conn {
	txByCtx := ctx.Value(types.TxKey)
	if txByCtx == nil {
		return t.db
	}
	tx, ok := txByCtx.(*sqlx.Tx)
	if ok {
		return tx
	}
	return t.db
}

type conn interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
