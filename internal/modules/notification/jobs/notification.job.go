package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"movie-app-go/internal/constants"
	"movie-app-go/internal/modules/notification/repositories"
	"movie-app-go/internal/modules/notification/services"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
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
	NotificationService *services.NotificationService
}

func NewNotificationJobHandler(db *gorm.DB) *NotificationJobHandler {
	notificationRepo := repositories.NewNotificationRepository(db)
	notificationService := services.NewNotificationService(notificationRepo)

	return &NotificationJobHandler{
		NotificationService: notificationService,
	}
}

func (h *NotificationJobHandler) HandleMovieReminder(ctx context.Context, t *asynq.Task) error {
	var payload MovieReminderPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal movie reminder payload: %w", err)
	}

	log.Printf("Processing movie reminder for user %d, movie: %s", payload.UserID, payload.MovieTitle)

	err := h.NotificationService.CreateMovieReminderNotification(
		payload.UserID,
		payload.MovieTitle,
		payload.ScheduleTime,
		payload.MovieID,
	)

	if err != nil {
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

	err := h.NotificationService.CreateBulkNotifications(
		payload.UserIDs,
		"New Promo Available!",
		fmt.Sprintf("Check out our new promo: %s. Use code: %s", payload.PromoName, payload.PromoCode),
		constants.NotificationTypePromoAvailable,
		map[string]interface{}{
			"promo_id":   payload.PromoID,
			"promo_name": payload.PromoName,
			"promo_code": payload.PromoCode,
		},
	)

	if err != nil {
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

	err := h.NotificationService.CreateBookingConfirmationNotification(
		payload.UserID,
		payload.TransactionID,
		payload.MovieTitle,
		payload.TotalAmount,
	)

	if err != nil {
		log.Printf("Failed to create booking confirmation notification: %v", err)
		return err
	}

	log.Printf("Booking confirmation notification sent to user %d", payload.UserID)
	return nil
}
