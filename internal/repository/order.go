package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/types"
)

type orderPostgres struct {
	tm transactionManager
}

func newOrderPostgres(tm transactionManager) *orderPostgres {
	return &orderPostgres{
		tm: tm,
	}
}
func (o *orderPostgres) FindUserIDByOrder(ctx context.Context, number string) (int64, error) {
	var userID int64
	query := "select user_id from orders where number=$1"
	if err := o.tm.getConn(ctx).GetContext(ctx, &userID, query, number); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return userID, nil
}

func (o *orderPostgres) CreateOrder(ctx context.Context, number string) error {
	userID := ctx.Value(types.UserIDKey).(int64)

	if _, err := o.tm.getConn(ctx).ExecContext(ctx, "insert into orders (user_id, number) values($1,$2)", userID, number); err != nil {
		return err
	}
	return nil
}

func (o *orderPostgres) FindOrdersByUserID(ctx context.Context, userID int64) ([]entities.Order, error) {
	var orders []entities.Order
	query := "select number, status , upload_at, accrual from orders where user_id = $1 order by upload_at"
	if err := o.tm.getConn(ctx).SelectContext(ctx, &orders, query, userID); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderPostgres) GetOrdersWithoutFinalStatuses(ctx context.Context) ([]entities.Order, error) {
	var orders []entities.Order
	query := `select number, status , upload_at, accrual  FROM
                                       orders WHERE status NOT IN ('INVALID', 'PROCESSED') ORDER BY upload_at`
	if err := o.tm.getConn(ctx).SelectContext(ctx, &orders, query); err != nil {
		return nil, err
	}
	return orders, nil

}

func (o *orderPostgres) GetOrders(ctx context.Context) ([]entities.Order, error) {
	var orders []entities.Order
	query := `select number, status , upload_at , accrual FROM
                                       orders  ORDER BY upload_at`
	if err := o.tm.getConn(ctx).SelectContext(ctx, &orders, query); err != nil {
		return nil, err
	}
	return orders, nil

}
func (o *orderPostgres) UpdateOrder(ctx context.Context, order entities.OrderAccrual) error {
	query := "update orders set status=$1, accrual=$2 where number=$3"
	_, err := o.tm.getConn(ctx).ExecContext(ctx, query, order.Status, order.Accrual, order.Order)
	return err
}
