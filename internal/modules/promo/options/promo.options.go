package options

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PromoOptions struct {
	Page         int    `json:"page"`
	PerPage      int    `json:"per_page"`
	Search       string `json:"search"`
	IsActive     *bool  `json:"is_active"`
	DiscountType string `json:"discount_type"`
	MovieID      *uint  `json:"movie_id"`
	ValidOnly    bool   `json:"valid_only"`
}

func ParsePromoOptions(c *gin.Context) (PromoOptions, error) {
	opts := PromoOptions{
		Page:    1,
		PerPage: 10,
	}

	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && page > 0 {
		opts.Page = page
	}

	if perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10")); err == nil && perPage > 0 {
		opts.PerPage = perPage
	}

	opts.Search = c.Query("search")
	opts.DiscountType = c.Query("discount_type")

	if isActive := c.Query("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			opts.IsActive = &active
		}
	}

	if movieIDStr := c.Query("movie_id"); movieIDStr != "" {
		if movieID, err := strconv.ParseUint(movieIDStr, 10, 32); err == nil {
			movieIDUint := uint(movieID)
			opts.MovieID = &movieIDUint
		}
	}

	if validOnly, err := strconv.ParseBool(c.DefaultQuery("valid_only", "false")); err == nil {
		opts.ValidOnly = validOnly
	}

	return opts, nil
}
