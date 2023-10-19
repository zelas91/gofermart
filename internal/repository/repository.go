package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/zelas91/gofermart/internal/entities"
)

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: newAuthPostgres(db),
	}
}

type Authorization interface {
	CreateUser(ctx context.Context, login, password string) error
	GetUser(ctx context.Context, user *entities.User) (entities.User, error)
}
