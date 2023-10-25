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
	return 0, nil
}

func (o *orderPostgres) CreateOrder(ctx context.Context, number string) error {
	userID := ctx.Value(types.UserIDKey).(int64)

	if _, err := o.db.ExecContext(ctx, "insert into orders (user_id, number) values($1,$2)", userID, number); err != nil {
		return err
	}
	return nil
}
func (o *orderPostgres) GetOrders(ctx context.Context, userID int64) ([]entities.Order, error) {
	var orders []entities.Order
	rows, err := o.db.QueryxContext(ctx, "select number, status , upload_at from orders where user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer func(ctx context.Context) {
		if err = rows.Close(); err != nil {
			logger.GetLogger(ctx).Errorf("rows close err: %v", err)
		}
	}(ctx)

	for rows.Next() {
		var order entities.Order
		if err = rows.Scan(&order.Number, &order.Status, &order.UploadedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)

	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	logger.GetLogger(ctx).Infof("len %d [] - > %v", len(orders), orders)
	return orders, nil
}
