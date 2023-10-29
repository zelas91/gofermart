package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zelas91/gofermart/internal/entities"
)

type balancePostgres struct {
	db *sqlx.DB
}

func newBalancePostgres(db *sqlx.DB) *balancePostgres {
	return &balancePostgres{
		db: db,
	}
}

func (b *balancePostgres) GetBalance(ctx context.Context, userID int64) (entities.Balance, error) {
	var balance entities.Balance
	query := "select sum(accrual)  as current from orders  where user_id=$1;"
	err := b.db.GetContext(ctx, &balance, query, userID)
	return balance, err
}
