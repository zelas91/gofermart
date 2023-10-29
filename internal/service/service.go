package service

import (
	"context"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/repository"
)

type Service struct {
	Authorization
	Orders
	Balance
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(repos),
		Orders:        newOrderService(repos),
		Balance:       newBalanceService(repos),
	}
}

type Balance interface {
	GetBalance(ctx context.Context) (entities.Balance, error)
}
type Authorization interface {
	CreateUser(ctx context.Context, user *entities.User) error
	CreateToken(ctx context.Context, user *entities.User) (string, error)
	ParserToken(ctx context.Context, tokenString string) (*entities.User, error)
}
type Orders interface {
	ValidateNumber(number string) bool
	FindUserIDByOrder(ctx context.Context, number string) (int64, error)
	CreateOrder(ctx context.Context, number string) error
	FindOrdersByUser(ctx context.Context) ([]entities.Order, error)
	GetOrders(ctx context.Context) ([]entities.Order, error)
	GetOrdersWithoutFinalStatuses(ctx context.Context) ([]entities.Order, error)
	UpdateOrder(ctx context.Context, order entities.OrderAccrual) error
}
