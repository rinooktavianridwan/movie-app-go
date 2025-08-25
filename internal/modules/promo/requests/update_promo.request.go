package requests

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type UpdatePromoRequest struct {
	Name          *string    `json:"name" binding:"omitempty,max=255"`
	Description   *string    `json:"description"`
	DiscountType  *string    `json:"discount_type" binding:"omitempty,oneof=percentage fixed_amount"`
	DiscountValue *float64   `json:"discount_value" binding:"omitempty,gt=0"`
	MinTickets    *int       `json:"min_tickets" binding:"omitempty,gte=1"`
	MaxDiscount   *float64   `json:"max_discount" binding:"omitempty,gt=0"`
	UsageLimit    *int       `json:"usage_limit" binding:"omitempty,gt=0"`
	IsActive      *bool      `json:"is_active"`
	ValidFrom     *time.Time `json:"valid_from"`
	ValidUntil    *time.Time `json:"valid_until"`
	MovieIDs      []uint     `json:"movie_ids"`
}

func (r *UpdatePromoRequest) Validate() error {
	validate := validator.New()
	if err := validate.Struct(r); err != nil {
		return err
	}

	if r.ValidFrom != nil && r.ValidUntil != nil {
		if !r.ValidUntil.After(*r.ValidFrom) {
			return validator.ValidationErrors{}
		}
	}

	return nil
}
