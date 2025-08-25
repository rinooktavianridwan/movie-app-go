package requests

type ValidatePromoRequest struct {
	PromoCode   string  `json:"promo_code" binding:"required"`
	TotalAmount float64 `json:"total_amount" binding:"required,gt=0"`
	MovieIDs    []uint  `json:"movie_ids" binding:"required"`
	SeatNumbers []uint  `json:"seat_numbers" binding:"required"`
}
