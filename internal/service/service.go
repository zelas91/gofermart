package service

import (
	"context"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/repository"
)

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: newAuthService(repos.Authorization),
	}
}

type Authorization interface {
	CreateUser(ctx context.Context, user *entities.User) error
	CreateToken(ctx context.Context, user *entities.User) (string, error)
	ParserToken(ctx context.Context, tokenString string) (*entities.User, error)
}
