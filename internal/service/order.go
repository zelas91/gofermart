package service

import (
	"context"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/repository"
	"github.com/zelas91/gofermart/internal/types"
	"strconv"
)

type orderService struct {
	repo repository.Orders
}

func newOrderService(repo repository.Orders) *orderService {
	return &orderService{repo: repo}
}
func (o *orderService) FindUserIDByOrder(ctx context.Context, number string) (int64, error) {
	return o.repo.FindUserIDByOrder(ctx, number)
}

func (o *orderService) CreateOrder(ctx context.Context, number string) error {
	if err := o.repo.CreateOrder(ctx, number); err != nil {
		return err
	}
	return nil
}

func (o *orderService) FindOrdersByUser(ctx context.Context) ([]entities.Order, error) {
	return o.repo.FindOrdersByUserID(ctx, ctx.Value(types.UserIDKey).(int64))
}
func (o *orderService) ValidateNumber(number string) bool {
	sum := 0
	parity := len(number) % 2
	for i := 0; i < len(number); i += 1 {
		digit, _ := strconv.Atoi(string(number[i]))
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
func (o *orderService) GetOrders(ctx context.Context) ([]entities.Order, error) {
	return o.repo.GetOrders(ctx)
}

func (o *orderService) GetOrdersWithoutFinalStatuses(ctx context.Context) ([]entities.Order, error) {
	return o.repo.GetOrdersWithoutFinalStatuses(ctx)
}

func (o *orderService) UpdateOrder(ctx context.Context, order entities.OrderAccrual) error {
	return o.repo.UpdateOrder(ctx, order)
}
