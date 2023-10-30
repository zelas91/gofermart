package service

import (
	"context"
	"github.com/zelas91/gofermart/internal/entities"
	"github.com/zelas91/gofermart/internal/repository"
	"github.com/zelas91/gofermart/internal/types"
)

type balanceService struct {
	repo repository.Balance
}

func newBalanceService(repo repository.Balance) *balanceService {
	return &balanceService{repo: repo}
}

func (b *balanceService) GetBalance(ctx context.Context) (entities.Balance, error) {
	userID := ctx.Value(types.UserIDKey).(int64)
	return b.repo.GetBalance(ctx, userID)
}
func (b *balanceService) Withdraw(ctx context.Context, withdraw entities.Withdraw) error {
	userID := ctx.Value(types.UserIDKey).(int64)
	return b.repo.Withdraw(ctx, userID, withdraw)
}

func (b *balanceService) WithdrawInfo(ctx context.Context) ([]entities.WithdrawInfo, error) {
	userID := ctx.Value(types.UserIDKey).(int64)
	return b.repo.WithdrawInfo(ctx, userID)
}
