package requests

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type CreatePromoRequest struct {
	Name          string    `json:"name" binding:"required,max=255"`
	Code          string    `json:"code" binding:"required,max=100"`
	Description   string    `json:"description"`
	DiscountType  string    `json:"discount_type" binding:"required,oneof=percentage fixed_amount"`
	DiscountValue float64   `json:"discount_value" binding:"required,gt=0"`
	MinTickets    int       `json:"min_tickets" binding:"gte=1"`
	MaxDiscount   *float64  `json:"max_discount" binding:"omitempty,gt=0"`
	UsageLimit    *int      `json:"usage_limit" binding:"omitempty,gt=0"`
	IsActive      bool      `json:"is_active"`
	ValidFrom     time.Time `json:"valid_from" binding:"required"`
	ValidUntil    time.Time `json:"valid_until" binding:"required"`
	MovieIDs      []uint    `json:"movie_ids"`
}

func (r *CreatePromoRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	if !r.ValidUntil.After(r.ValidFrom) {
		return validator.ValidationErrors{}
	}

	return nil
}
