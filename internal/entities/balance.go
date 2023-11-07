package entities

import "github.com/zelas91/gofermart/internal/types"

type Balance struct {
	Current   float64 `json:"current" db:"current"`
	Withdrawn float64 `json:"withdrawn" db:"withdrawn"`
}

type Withdraw struct {
	Sum   float64 `json:"sum" validate:"required"`
	Order string  `json:"order" validate:"required"`
}
type WithdrawInfo struct {
	Order       string            `json:"order" db:"number"`
	Accrual     float64           `json:"sum" db:"accrual"`
	ProcessedAt types.RFC3339Time `json:"processed_at" db:"upload_at"`
}
