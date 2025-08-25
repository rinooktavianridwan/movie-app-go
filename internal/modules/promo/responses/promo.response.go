package responses

import (
	"movie-app-go/internal/models"
	"time"
)

type PromoResponse struct {
	ID            uint                 `json:"id"`
	Name          string               `json:"name"`
	Code          string               `json:"code"`
	Description   string               `json:"description"`
	DiscountType  string               `json:"discount_type"`
	DiscountValue float64              `json:"discount_value"`
	MinTickets    int                  `json:"min_tickets"`
	MaxDiscount   *float64             `json:"max_discount"`
	UsageLimit    *int                 `json:"usage_limit"`
	UsageCount    int                  `json:"usage_count"`
	IsActive      bool                 `json:"is_active"`
	ValidFrom     time.Time            `json:"valid_from"`
	ValidUntil    time.Time            `json:"valid_until"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
	Movies        []PromoMovieResponse `json:"movies,omitempty"`
}

type PromoMovieResponse struct {
	MovieID    uint   `json:"movie_id"`
	MovieTitle string `json:"movie_title"`
}

type PaginatedPromoResponse struct {
	Page      int             `json:"page"`
	PerPage   int             `json:"per_page"`
	Total     int64           `json:"total"`
	TotalPage int             `json:"total_pages"`
	Data      []PromoResponse `json:"data"`
}

type PromoValidationResponse struct {
	IsValid        bool    `json:"is_valid"`
	DiscountAmount float64 `json:"discount_amount"`
	FinalAmount    float64 `json:"final_amount"`
	Message        string  `json:"message"`
}

func ToPromoResponse(promo models.Promo) PromoResponse {
	response := PromoResponse{
		ID:            promo.ID,
		Name:          promo.Name,
		Code:          promo.Code,
		Description:   promo.Description,
		DiscountType:  promo.DiscountType,
		DiscountValue: promo.DiscountValue,
		MinTickets:    promo.MinTickets,
		MaxDiscount:   promo.MaxDiscount,
		UsageLimit:    promo.UsageLimit,
		UsageCount:    promo.UsageCount,
		IsActive:      promo.IsActive,
		ValidFrom:     promo.ValidFrom,
		ValidUntil:    promo.ValidUntil,
		CreatedAt:     promo.CreatedAt,
		UpdatedAt:     promo.UpdatedAt,
	}

	if len(promo.PromoMovies) > 0 {
		for _, pm := range promo.PromoMovies {
			response.Movies = append(response.Movies, PromoMovieResponse{
				MovieID:    pm.MovieID,
				MovieTitle: pm.Movie.Title,
			})
		}
	}

	return response
}

func ToPromoResponses(promos []models.Promo) []PromoResponse {
	responses := make([]PromoResponse, len(promos))
	for i, promo := range promos {
		responses[i] = ToPromoResponse(promo)
	}
	return responses
}
