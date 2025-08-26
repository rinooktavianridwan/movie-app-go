package options

import (
    "fmt"
    "strconv"
    "github.com/gin-gonic/gin"
)

type NotificationOptions struct {
    Page     int    `json:"page"`
    PerPage  int    `json:"per_page"`
    Type     string `json:"type,omitempty"`
    IsRead   *bool  `json:"is_read,omitempty"`
    UserID   uint   `json:"user_id,omitempty"`
}

func ParseNotificationOptions(ctx *gin.Context) (NotificationOptions, error) {
    options := NotificationOptions{
        Page:    1,
        PerPage: 10,
    }

    if page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")); err == nil && page > 0 {
        options.Page = page
    }
    if perPage, err := strconv.Atoi(ctx.DefaultQuery("per_page", "10")); err == nil && perPage > 0 {
        options.PerPage = perPage
    }

    if notificationType := ctx.Query("type"); notificationType != "" {
        validTypes := []string{"movie_reminder", "promo_available", "booking_confirmation", "system"}
        isValid := false
        for _, validType := range validTypes {
            if notificationType == validType {
                isValid = true
                break
            }
        }
        if !isValid {
            return options, fmt.Errorf("invalid notification type")
        }
        options.Type = notificationType
    }

    if isReadStr := ctx.Query("is_read"); isReadStr != "" {
        if isRead, err := strconv.ParseBool(isReadStr); err == nil {
            options.IsRead = &isRead
        } else {
            return options, fmt.Errorf("invalid is_read format, use true or false")
        }
    }

    return options, nil
}
