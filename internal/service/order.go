package service

import (
	"context"
	"github.com/zelas91/gofermart/internal/repository"
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
	return o.repo.CreateOrder(ctx, number)
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
