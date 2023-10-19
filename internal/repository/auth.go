package repository

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/zelas91/gofermart/internal/entities"
)

var (
	pgErr        *pq.Error
	ErrDuplicate = errors.New("login is already taken")
)

type AuthPostgres struct {
	db *sqlx.DB
}

func newAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(ctx context.Context, login, password string) error {
	if _, err := a.db.ExecContext(ctx,
		"INSERT INTO USERS (login, password) values($1, $2)", login, password); err != nil {
		if errors.As(err, &pgErr) {
			if pgError, ok := err.(*pq.Error); ok && pgError.Code == "23505" {
				return ErrDuplicate
			}

		}
		return err
	}
	return nil
}
func (a *AuthPostgres) GetUser(ctx context.Context, authUser *entities.User) (entities.User, error) {
	user := entities.User{}

	if err := a.db.GetContext(ctx, &user,
		"SELECT * FROM users WHERE login=$1", authUser.Login); err != nil {
		return user, err
	}
	return user, nil
}
