package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/logger"
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
	var userID int64
	query := "select user_id from orders where number=$1"
	if err := o.db.GetContext(ctx, &userID, query, number); err != nil {
		return 0, err
	}
	return userID, nil
}

func (o *orderPostgres) CreateOrder(ctx context.Context, number string) error {
	userID := ctx.Value(types.UserIDKey).(int64)
	var userIDReturn int64
	query := `with e as (
		insert  into  orders (user_id, number) VALUES ($1, $2) on conflict ("number") do nothing returning user_id
			) select * from e union select user_id from orders where number=$2`
	if err := o.db.GetContext(ctx, &userIDReturn, query, userID, number); err != nil {
		return err
	}
	logger.GetLogger(ctx).Info("!!!!!!! ", userID, " : ", userIDReturn, " NUMBER: ", number)
	return nil
}

func (o *orderPostgres) FindOrdersByUserID(ctx context.Context, userID int64) ([]entities.Order, error) {
	var orders []entities.Order
	query := "select number, status , upload_at from orders where user_id = $1 order by upload_at"
	if err := o.db.SelectContext(ctx, &orders, query, userID); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderPostgres) GetOrdersWithoutFinalStatuses(ctx context.Context) ([]entities.Order, error) {
	var orders []entities.Order
	query := `select number, status , upload_at  FROM
                                       orders WHERE status NOT IN ('INVALID', 'PROCESSED') ORDER BY upload_at`
	if err := o.db.SelectContext(ctx, &orders, query); err != nil {
		return nil, err
	}
	return orders, nil

}

func (o *orderPostgres) GetOrders(ctx context.Context) ([]entities.Order, error) {
	var orders []entities.Order
	query := `select number, status , upload_at  FROM
                                       orders  ORDER BY upload_at`
	if err := o.db.SelectContext(ctx, &orders, query); err != nil {
		return nil, err
	}
	return orders, nil

}
