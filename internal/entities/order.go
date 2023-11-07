package entities

import "github.com/zelas91/gofermart/internal/types"

type Order struct {
	Number     string            `json:"number" db:"number"`
	Status     string            `json:"status" db:"status"`
	Accrual    float64           `json:"accrual,omitempty" db:"accrual"`
	UploadedAt types.RFC3339Time `json:"uploaded_at" db:"upload_at"`
}

type OrderAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
