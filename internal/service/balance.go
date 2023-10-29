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
