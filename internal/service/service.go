package service

import (
	"context"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/repository"
)

type Service struct {
	Authorization
	Orders
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(repos),
		Orders:        newOrderService(repos),
	}
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
	GetOrders(ctx context.Context) ([]entities.Order, error)
}
