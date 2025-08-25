package models

import (
	"time"

	"gorm.io/gorm"
)

type Promo struct {
	gorm.Model
	Name          string    `json:"name" gorm:"not null"`
	Code          string    `json:"code" gorm:"uniqueIndex;not null"`
	Description   string    `json:"description"`
	DiscountType  string    `json:"discount_type" gorm:"not null"`
	DiscountValue float64   `json:"discount_value" gorm:"not null"`
	MinTickets    int       `json:"min_tickets" gorm:"default:1"`
	MaxDiscount   *float64  `json:"max_discount"`
	UsageLimit    *int      `json:"usage_limit"`
	UsageCount    int       `json:"usage_count" gorm:"default:0"`
	IsActive      bool      `json:"is_active" gorm:"default:true"`
	ValidFrom     time.Time `json:"valid_from"`
	ValidUntil    time.Time `json:"valid_until"`

	PromoMovies []PromoMovie `json:"promo_movies,omitempty" gorm:"foreignKey:PromoID"`
	PromoUsage  []PromoUsage `json:"promo_usage,omitempty" gorm:"foreignKey:PromoID"`
}

type PromoMovie struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PromoID   uint      `json:"promo_id" gorm:"not null"`
	MovieID   uint      `json:"movie_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`

	Promo Promo `json:"promo,omitempty" gorm:"foreignKey:PromoID"`
	Movie Movie `json:"movie,omitempty" gorm:"foreignKey:MovieID"`
}

type PromoUsage struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	PromoID        uint      `json:"promo_id" gorm:"not null"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	TransactionID  uint      `json:"transaction_id" gorm:"not null;uniqueIndex"`
	DiscountAmount float64   `json:"discount_amount" gorm:"not null"`
	UsedAt         time.Time `json:"used_at" gorm:"default:CURRENT_TIMESTAMP"`

	Promo       Promo       `json:"promo,omitempty" gorm:"foreignKey:PromoID"`
	User        User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Transaction Transaction `json:"transaction,omitempty" gorm:"foreignKey:TransactionID"`
}
