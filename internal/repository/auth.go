package repository

import (
	"context"
	"errors"
	"github.com/lib/pq"
	"github.com/zelas91/gofermart/internal/entities"
)

var (
	pgError      *pq.Error
	ErrDuplicate = errors.New("login is already taken")
)

type authPostgres struct {
	tm transactionManager
}

func newAuthPostgres(tm transactionManager) *authPostgres {
	return &authPostgres{tm: tm}
}

func (a *authPostgres) CreateUser(ctx context.Context, login, password string) error {
	if _, err := a.tm.getConn(ctx).ExecContext(ctx,
		"INSERT INTO USERS (login, password) values($1, $2)", login, password); err != nil {

		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return ErrDuplicate
		}

		return err
	}
	return nil
}
func (a *authPostgres) GetUser(ctx context.Context, authUser *entities.User) (entities.User, error) {
	user := entities.User{}

	if err := a.tm.getConn(ctx).GetContext(ctx, &user,
		"SELECT * FROM users WHERE login=$1", authUser.Login); err != nil {
		return user, err
	}
	return user, nil
}
