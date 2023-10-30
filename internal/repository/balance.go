package repository

import (
	"context"
	"errors"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/types"
)

var (
	ErrNotEnoughFunds = errors.New("not enough funds")
)

type balancePostgres struct {
	tm     transactionManager
	orders Orders
}

func newBalancePostgres(tm transactionManager, orders Orders) *balancePostgres {
	return &balancePostgres{
		tm:     tm,
		orders: orders,
	}
}

func (b *balancePostgres) GetBalance(ctx context.Context, userID int64) (entities.Balance, error) {
	var balance entities.Balance
	query := `select sum(o.accrual)  as current ,
       sum(case when o.accrual  < 0 then -1 * o.accrual  else 0 end) as withdrawn from orders o where user_id=$1;`
	err := b.tm.getConn(ctx).GetContext(ctx, &balance, query, userID)
	return balance, err
}

func (b *balancePostgres) Withdraw(ctx context.Context, userID int64, withdraw entities.Withdraw) error {
	return b.tm.do(ctx, func(ctx context.Context) error {
		usID, err := b.orders.FindUserIDByOrder(ctx, withdraw.Order)
		if err != nil {
			return err
		}
		if usID == 0 {
			balance, err := b.GetBalance(ctx, userID)
			if err != nil {
				return err
			}

			if balance.Current < 0 || balance.Current < withdraw.Sum {
				return ErrNotEnoughFunds
			}
			query := `insert into orders (number,  accrual, user_id, status) values ($1,$2,$3,$4)`
			if _, err = b.tm.getConn(ctx).ExecContext(ctx, query, withdraw.Order, -1*withdraw.Sum, userID, types.PROCESSED); err != nil {
				return err
			}
		} else if usID != userID {
			return errors.New("the order belongs to another user")
		}

		return nil
	})
}
func (b *balancePostgres) WithdrawInfo(ctx context.Context, userID int64) ([]entities.WithdrawInfo, error) {
	var withdraws []entities.WithdrawInfo
	query := `select number, (-1 *accrual) as accrual, upload_at from orders where user_id=$1 and accrual <0`
	if err := b.tm.getConn(ctx).SelectContext(ctx, &withdraws, query, userID); err != nil {
		return nil, err
	}
	return withdraws, nil

}
