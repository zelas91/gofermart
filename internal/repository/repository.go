package repository

import (
	"context"
	"github.com/zelas91/gofermart/internal/entities"
)

type Authorization interface {
	CreateUser(cxt context.Context, user entities.User) error
	GetUser(ctx context.Context, user entities.User) (entities.User, error)
}
