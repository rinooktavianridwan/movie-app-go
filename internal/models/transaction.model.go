package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID        uint    `json:"user_id"`
	TotalAmount   float64 `json:"total_amount"`
	PaymentMethod string  `json:"payment_method"`
	PaymentStatus string  `json:"payment_status"`
}
