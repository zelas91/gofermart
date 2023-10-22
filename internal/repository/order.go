package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zelas91/gofermart/internal/types"
)

type orderPostgres struct {
	db *sqlx.DB
}

func newOrderPostgres(db *sqlx.DB) *orderPostgres {
	return &orderPostgres{
		db: db,
	}
}
func (o *orderPostgres) FindUserIDByOrder(ctx context.Context, number string) (int64, error) {
	return 0, nil
}

func (o *orderPostgres) CreateOrder(ctx context.Context, number string) error {
	userID := ctx.Value(types.UserIDKey).(int64)

	if _, err := o.db.ExecContext(ctx, "insert into orders (user_id, number) values($1,$2)", userID, number); err != nil {
		return err
	}
	return nil
}
