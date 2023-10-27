package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zelas91/gofermart/internal/entities"
)

type Repository struct {
	Authorization
	Orders
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthPostgres(db),
		Orders:        newOrderPostgres(db),
	}
}

//go:generate mockgen -package mocks -destination=./mocks/mock_repository.go -source=repository.go -package=mock
type Authorization interface {
	CreateUser(ctx context.Context, login, password string) error
	GetUser(ctx context.Context, user *entities.User) (entities.User, error)
}

type Orders interface {
	FindUserIDByOrder(ctx context.Context, number string) (int64, error)
	CreateOrder(ctx context.Context, number string) error
	FindOrdersByUserID(ctx context.Context, userID int64) ([]entities.Order, error)
	GetOrders(ctx context.Context) ([]entities.Order, error)
	GetOrdersWithoutFinalStatuses(ctx context.Context) ([]entities.Order, error)
}
