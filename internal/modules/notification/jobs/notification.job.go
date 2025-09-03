package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const (
	TypeMovieReminder     = "notification:movie_reminder"
	TypePromoNotification = "notification:promo_notification"
	TypeBookingConfirm    = "notification:booking_confirmation"
)

type MovieReminderPayload struct {
	UserID       uint   `json:"user_id"`
	MovieID      uint   `json:"movie_id"`
	MovieTitle   string `json:"movie_title"`
	ScheduleTime string `json:"schedule_time"`
	ScheduleID   uint   `json:"schedule_id"`
}

type PromoNotificationPayload struct {
	UserIDs   []uint `json:"user_ids"`
	PromoID   uint   `json:"promo_id"`
	PromoName string `json:"promo_name"`
	PromoCode string `json:"promo_code"`
	MovieIDs  []uint `json:"movie_ids,omitempty"`
}

type BookingConfirmationPayload struct {
	UserID        uint    `json:"user_id"`
	TransactionID uint    `json:"transaction_id"`
	MovieTitle    string  `json:"movie_title"`
	TotalAmount   float64 `json:"total_amount"`
	ScheduleTime  string  `json:"schedule_time"`
}

type NotificationJobHandler struct {
	DB *gorm.DB
}

func NewNotificationJobHandler(db *gorm.DB) *NotificationJobHandler {
	return &NotificationJobHandler{
		DB: db,
	}
}

func (h *NotificationJobHandler) HandleMovieReminder(ctx context.Context, t *asynq.Task) error {
	var payload MovieReminderPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal movie reminder payload: %w", err)
	}

	log.Printf("Processing movie reminder for user %d, movie: %s", payload.UserID, payload.MovieTitle)

	notification := models.Notification{
		UserID:  payload.UserID,
		Title:   "Movie Reminder",
		Message: fmt.Sprintf("Your movie '%s' is starting soon at %s", payload.MovieTitle, payload.ScheduleTime),
		Type:    constants.NotificationTypeMovieReminder,
		IsRead:  false,
		Data: models.NotificationData{
			"movie_id":      payload.MovieID,
			"movie_title":   payload.MovieTitle,
			"schedule_time": payload.ScheduleTime,
		},
	}

	if err := h.DB.Create(&notification).Error; err != nil {
		log.Printf("Failed to create movie reminder notification: %v", err)
		return err
	}

	log.Printf("Movie reminder notification sent to user %d", payload.UserID)
	return nil
}

func (h *NotificationJobHandler) HandlePromoNotification(ctx context.Context, t *asynq.Task) error {
	var payload PromoNotificationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal promo notification payload: %w", err)
	}

	log.Printf("Processing promo notification for %d users, promo: %s", len(payload.UserIDs), payload.PromoName)

	notifications := make([]models.Notification, len(payload.UserIDs))
	for i, userID := range payload.UserIDs {
		notifications[i] = models.Notification{
			UserID:  userID,
			Title:   "New Promo Available!",
			Message: fmt.Sprintf("Check out our new promo: %s. Use code: %s", payload.PromoName, payload.PromoCode),
			Type:    constants.NotificationTypePromoAvailable,
			IsRead:  false,
			Data: models.NotificationData{
				"promo_id":   payload.PromoID,
				"promo_name": payload.PromoName,
				"promo_code": payload.PromoCode,
			},
		}
	}

	if err := h.DB.Create(&notifications).Error; err != nil {
		log.Printf("Failed to create promo notifications: %v", err)
		return err
	}

	log.Printf("Promo notification sent to %d users", len(payload.UserIDs))
	return nil
}

func (h *NotificationJobHandler) HandleBookingConfirmation(ctx context.Context, t *asynq.Task) error {
	var payload BookingConfirmationPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal booking confirmation payload: %w", err)
	}

	log.Printf("Processing booking confirmation for user %d, transaction %d", payload.UserID, payload.TransactionID)

	notification := models.Notification{
		UserID:  payload.UserID,
		Title:   "Booking Confirmed",
		Message: fmt.Sprintf("Your booking for '%s' has been confirmed. Total: Rp %.0f", payload.MovieTitle, payload.TotalAmount),
		Type:    constants.NotificationTypeBookingConfirm,
		IsRead:  false,
		Data: models.NotificationData{
			"transaction_id": payload.TransactionID,
			"movie_title":    payload.MovieTitle,
			"total_amount":   payload.TotalAmount,
		},
	}

	if err := h.DB.Create(&notification).Error; err != nil {
		log.Printf("Failed to create booking confirmation notification: %v", err)
		return err
	}

	log.Printf("Booking confirmation notification sent to user %d", payload.UserID)
	return nil
}
