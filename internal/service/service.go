package service

import (
	"context"
	"github.com/zelas91/gofermart/internal/entities"
)

type Authorization interface {
	CreateUser(cxt context.Context, user entities.User) error
	CreateToken(ctx context.Context, user entities.User) (string, error)
}
