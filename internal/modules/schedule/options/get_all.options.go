package options

import (
    "fmt"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
)

type GetAllScheduleOptions struct {
    Page       int        `json:"page"`
    PerPage    int        `json:"per_page"`
    StudioID   *uint      `json:"studio_id,omitempty"`
    MovieTitle string     `json:"movie_title,omitempty"`
    DateFrom   *time.Time `json:"date_from,omitempty"`
    DateTo     *time.Time `json:"date_to,omitempty"`
}



// ParseScheduleOptions - Parse query parameters ke GetAllScheduleOptions
func ParseScheduleOptions(ctx *gin.Context) (*GetAllScheduleOptions, error) {
    options := &GetAllScheduleOptions{
        Page:    1,
        PerPage: 10,
    }

    // Parse page dan per_page
    if page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")); err == nil && page > 0 {
        options.Page = page
    }
    if perPage, err := strconv.Atoi(ctx.DefaultQuery("per_page", "10")); err == nil && perPage > 0 {
        options.PerPage = perPage
    }

    // Parse studio_id
    if studioIDStr := ctx.Query("studio_id"); studioIDStr != "" {
        if studioID, err := strconv.Atoi(studioIDStr); err == nil {
            studioIDUint := uint(studioID)
            options.StudioID = &studioIDUint
        } else {
            return nil, fmt.Errorf("invalid studio_id format")
        }
    }

    // Parse movie_title
    options.MovieTitle = ctx.Query("movie_title")

    // Parse date_from
    if dateFromStr := ctx.Query("date_from"); dateFromStr != "" {
        if dateFrom, err := time.Parse("2006-01-02", dateFromStr); err == nil {
            options.DateFrom = &dateFrom
        } else {
            return nil, fmt.Errorf("invalid date_from format, use YYYY-MM-DD")
        }
    }

    // Parse date_to
    if dateToStr := ctx.Query("date_to"); dateToStr != "" {
        if dateTo, err := time.Parse("2006-01-02", dateToStr); err == nil {
            options.DateTo = &dateTo
        } else {
            return nil, fmt.Errorf("invalid date_to format, use YYYY-MM-DD")
        }
    }

    if options.DateFrom != nil && options.DateTo != nil {
        if options.DateFrom.After(*options.DateTo) {
            return nil, fmt.Errorf("date_from cannot be after date_to")
        }
    }

    return options, nil
}