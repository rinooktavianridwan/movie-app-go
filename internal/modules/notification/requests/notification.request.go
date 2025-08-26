package requests

type CreateNotificationRequest struct {
    UserID  uint                   `json:"user_id" binding:"required"`
    Title   string                 `json:"title" binding:"required,max=255"`
    Message string                 `json:"message" binding:"required"`
    Type    string                 `json:"type" binding:"required,oneof=movie_reminder promo_available booking_confirmation system"`
    Data    map[string]interface{} `json:"data,omitempty"`
}

type UpdateNotificationRequest struct {
	IsRead *bool `json:"is_read"`
}

type BulkMarkReadRequest struct {
	NotificationIDs []uint `json:"notification_ids" binding:"required"`
}

type CreateSystemNotificationRequest struct {
    Title   string                 `json:"title" binding:"required,max=255"`
    Message string                 `json:"message" binding:"required"`
    Data    map[string]interface{} `json:"data,omitempty"`
}
