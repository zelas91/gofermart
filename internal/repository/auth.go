package repository

import (
	"context"
	"errors"
	"github.com/lib/pq"
	"github.com/zelas91/gofermart/internal/entities"
	errorService "github.com/zelas91/gofermart/internal/error"
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
		if errPg := new(pq.PGError); errors.As(err, errPg) {
			return errorService.ErrDuplicate
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
