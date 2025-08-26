package responses

import (
    "movie-app-go/internal/models"
)

type NotificationResponse struct {
    ID        uint                   `json:"id"`
    UserID    uint                   `json:"user_id"`
    Title     string                 `json:"title"`
    Message   string                 `json:"message"`
    Type      string                 `json:"type"`
    IsRead    bool                   `json:"is_read"`
    Data      map[string]interface{} `json:"data,omitempty"`
}

type PaginatedNotificationResponse struct {
    Page      int                    `json:"page"`
    PerPage   int                    `json:"per_page"`
    Total     int64                  `json:"total"`
    TotalPage int                    `json:"total_pages"`
    Data      []NotificationResponse `json:"data"`
}

type NotificationStatsResponse struct {
    TotalCount  int64 `json:"total_count"`
    UnreadCount int64 `json:"unread_count"`
    ReadCount   int64 `json:"read_count"`
}

func ToNotificationResponse(notification models.Notification) NotificationResponse {
    return NotificationResponse{
        ID:        notification.ID,
        UserID:    notification.UserID,
        Title:     notification.Title,
        Message:   notification.Message,
        Type:      notification.Type,
        IsRead:    notification.IsRead,
        Data:      notification.Data,
    }
}

func ToNotificationResponses(notifications []models.Notification) []NotificationResponse {
    responses := make([]NotificationResponse, len(notifications))
    for i, notification := range notifications {
        responses[i] = ToNotificationResponse(notification)
    }
    return responses
}
