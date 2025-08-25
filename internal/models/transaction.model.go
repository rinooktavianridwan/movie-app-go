package models

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID         uint     `json:"user_id" gorm:"not null"`
	TotalAmount    float64  `json:"total_amount" gorm:"not null"`
	OriginalAmount *float64 `json:"original_amount"`
	DiscountAmount float64  `json:"discount_amount" gorm:"default:0"`
	PaymentMethod  string   `json:"payment_method" gorm:"not null"`
	PaymentStatus  string   `json:"payment_status" gorm:"not null;default:'pending'"`
	PromoID        *uint    `json:"promo_id"`

	User    User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Promo   *Promo   `json:"promo,omitempty" gorm:"foreignKey:PromoID"`
	Tickets []Ticket `json:"tickets,omitempty" gorm:"foreignKey:TransactionID"`
}
