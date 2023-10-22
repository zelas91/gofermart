package entities

import "github.com/zelas91/gofermart/internal/types"

type Order struct {
	Number     string            `json:"number" db:"number"`
	Status     types.OrderStatus `json:"status" db:"status"`
	Accrual    float64           `json:"accrual,omitempty"`
	UploadedAt types.RFC3339Time `json:"uploaded_at" db:"upload_at"`
}
