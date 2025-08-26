package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/notification/services"

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
	DB                  *gorm.DB
	NotificationService *services.NotificationService
}

func NewNotificationJobHandler(db *gorm.DB) *NotificationJobHandler {
	notificationService := services.NewNotificationService(db)
	return &NotificationJobHandler{
		DB:                  db,
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

	successCount := 0
	for _, userID := range payload.UserIDs {
		err := h.NotificationService.CreatePromoNotification(
			userID,
			payload.PromoName,
			payload.PromoCode,
			payload.PromoID,
		)

		if err != nil {
			log.Printf("Failed to send promo notification to user %d: %v", userID, err)
			continue
		}
		successCount++
	}

	log.Printf("Promo notification sent to %d/%d users", successCount, len(payload.UserIDs))
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

func (h *NotificationJobHandler) ScheduleMovieReminder(userID uint, movieID uint, movieTitle string, scheduleTime time.Time, scheduleID uint) error {
	payload := MovieReminderPayload{
		UserID:       userID,
		MovieID:      movieID,
		MovieTitle:   movieTitle,
		ScheduleTime: scheduleTime.Format("2006-01-02 15:04:05"),
		ScheduleID:   scheduleID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	reminderTime := scheduleTime.Add(-time.Hour)
	log.Printf("Scheduled movie reminder for user %d at %v (payload: %s)", userID, reminderTime, string(payloadBytes))

	return nil
}

func (h *NotificationJobHandler) SchedulePromoNotification(userIDs []uint, promoID uint, promoName string, promoCode string, movieIDs []uint, scheduleTime time.Time) error {
	payload := PromoNotificationPayload{
		UserIDs:   userIDs,
		PromoID:   promoID,
		PromoName: promoName,
		PromoCode: promoCode,
		MovieIDs:  movieIDs,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Printf("Scheduled promo notification for %d users at %v (payload: %s)", len(userIDs), scheduleTime, string(payloadBytes))

	return nil
}

func (h *NotificationJobHandler) SendBookingConfirmation(userID uint, transactionID uint, movieTitle string, totalAmount float64, scheduleTime string) error {
	payload := BookingConfirmationPayload{
		UserID:        userID,
		TransactionID: transactionID,
		MovieTitle:    movieTitle,
		TotalAmount:   totalAmount,
		ScheduleTime:  scheduleTime,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	log.Printf("Queued booking confirmation for user %d, transaction %d (payload: %s)", userID, transactionID, string(payloadBytes))

	return nil
}

func (h *NotificationJobHandler) GetActiveUsersForPromoNotification(movieIDs []uint) ([]uint, error) {
	var userIDs []uint

	query := h.DB.Model(&models.User{}).
		Where("is_admin = ?", false).
		Select("id")

	if len(movieIDs) > 0 {
		query = query.Joins("JOIN transactions ON transactions.user_id = users.id").
			Joins("JOIN schedules ON schedules.id = transactions.schedule_id").
			Where("schedules.movie_id IN ?", movieIDs).
			Distinct()
	}

	if err := query.Pluck("id", &userIDs).Error; err != nil {
		return nil, err
	}

	return userIDs, nil
}
