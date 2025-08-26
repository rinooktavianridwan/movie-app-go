package services

import (
	"fmt"
	"log"
	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/notification/options"
	"movie-app-go/internal/modules/notification/requests"
	"movie-app-go/internal/modules/notification/responses"
	"movie-app-go/internal/repository"

	"gorm.io/gorm"
)

type NotificationService struct {
	DB *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{
		DB: db,
	}
}

func (s *NotificationService) CreateNotification(req *requests.CreateNotificationRequest) (*models.Notification, error) {
	notification := models.Notification{
		UserID:  req.UserID,
		Title:   req.Title,
		Message: req.Message,
		Type:    req.Type,
		IsRead:  false,
		Data:    req.Data,
	}

	if err := s.DB.Create(&notification).Error; err != nil {
		return nil, err
	}

	return &notification, nil
}

func (s *NotificationService) GetUserNotifications(userID uint, opts options.NotificationOptions) (repository.PaginationResult[models.Notification], error) {
	var notifications []models.Notification
	var total int64

	query := s.DB.Model(&models.Notification{}).Where("user_id = ?", userID)

	if opts.Type != "" {
		query = query.Where("type = ?", opts.Type)
	}

	if opts.IsRead != nil {
		query = query.Where("is_read = ?", *opts.IsRead)
	}

	if err := query.Count(&total).Error; err != nil {
		return repository.PaginationResult[models.Notification]{}, err
	}

	offset := (opts.Page - 1) * opts.PerPage
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(opts.PerPage).
		Find(&notifications).Error; err != nil {
		return repository.PaginationResult[models.Notification]{}, err
	}

	totalPages := int(total) / opts.PerPage
	if int(total)%opts.PerPage != 0 {
		totalPages++
	}

	return repository.PaginationResult[models.Notification]{
		Data:       notifications,
		Total:      total,
		Page:       opts.Page,
		PerPage:    opts.PerPage,
		TotalPages: totalPages,
	}, nil
}

func (s *NotificationService) GetNotificationByID(id uint) (*models.Notification, error) {
	var notification models.Notification
	if err := s.DB.First(&notification, id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (s *NotificationService) MarkAsRead(id uint, userID uint) (*models.Notification, error) {
	var notification models.Notification

	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&notification).Error; err != nil {
		return nil, fmt.Errorf("notification not found")
	}

	notification.IsRead = true
	if err := s.DB.Save(&notification).Error; err != nil {
		return nil, err
	}

	return &notification, nil
}

func (s *NotificationService) MarkAsUnread(id uint, userID uint) (*models.Notification, error) {
	var notification models.Notification

	if err := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&notification).Error; err != nil {
		return nil, fmt.Errorf("notification not found")
	}

	notification.IsRead = false
	if err := s.DB.Save(&notification).Error; err != nil {
		return nil, err
	}

	return &notification, nil
}

func (s *NotificationService) BulkMarkAsRead(req *requests.BulkMarkReadRequest, userID uint) error {
	if len(req.NotificationIDs) == 0 {
		return fmt.Errorf("no notification IDs provided")
	}

	result := s.DB.Model(&models.Notification{}).
		Where("id IN ? AND user_id = ?", req.NotificationIDs, userID).
		Update("is_read", true)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no notifications found to update")
	}

	return nil
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}

func (s *NotificationService) GetNotificationStats(userID uint) (*responses.NotificationStatsResponse, error) {
	var totalCount, unreadCount int64

	if err := s.DB.Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if err := s.DB.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&unreadCount).Error; err != nil {
		return nil, err
	}

	return &responses.NotificationStatsResponse{
		TotalCount:  totalCount,
		UnreadCount: unreadCount,
		ReadCount:   totalCount - unreadCount,
	}, nil
}

func (s *NotificationService) DeleteNotification(id uint, userID uint) error {
	result := s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

func (s *NotificationService) CreateMovieReminderNotification(userID uint, movieTitle string, scheduleTime string, movieID uint) error {
	req := &requests.CreateNotificationRequest{
		UserID:  userID,
		Title:   "Movie Reminder",
		Message: fmt.Sprintf("Your movie '%s' is starting soon at %s", movieTitle, scheduleTime),
		Type:    constants.NotificationTypeMovieReminder,
		Data: map[string]interface{}{
			"movie_id":      movieID,
			"movie_title":   movieTitle,
			"schedule_time": scheduleTime,
		},
	}

	_, err := s.CreateNotification(req)
	return err
}

func (s *NotificationService) CreatePromoNotification(promoID uint, promoName, promoCode string) error {
    var userIDs []uint
    if err := s.DB.Model(&models.User{}).
        Where("is_admin = ?", false).
        Pluck("id", &userIDs).Error; err != nil {
        log.Printf("Failed to get users for promo notification: %v", err)
        return err
    }

    if len(userIDs) == 0 {
        log.Printf("No users found for promo notification")
        return fmt.Errorf("no users found")
    }

    return s.CreateBulkNotifications(
        userIDs,
        "New Promo Available!",
        fmt.Sprintf("Check out our new promo: %s. Use code: %s", promoName, promoCode),
        constants.NotificationTypePromoAvailable,
        map[string]interface{}{
            "promo_id":   promoID,
            "promo_name": promoName,
            "promo_code": promoCode,
        },
    )
}

func (s *NotificationService) CreateBookingConfirmationNotification(
	userID uint,
	transactionID uint,
	movieTitle string,
	totalAmount float64,
) error {
	req := &requests.CreateNotificationRequest{
		UserID:  userID,
		Title:   "Booking Confirmed",
		Message: fmt.Sprintf("Your booking for '%s' has been confirmed. Total: Rp %.0f", movieTitle, totalAmount),
		Type:    constants.NotificationTypeBookingConfirm,
		Data: map[string]interface{}{
			"transaction_id": transactionID,
			"movie_title":    movieTitle,
			"total_amount":   totalAmount,
		},
	}

	_, err := s.CreateNotification(req)
	return err
}

func (s *NotificationService) CreateBulkNotifications(
	userIDs []uint,
	title string,
	message string,
	notificationType string,
	data map[string]interface{},
) error {
	if len(userIDs) == 0 {
		return fmt.Errorf("no user IDs provided")
	}

	notifications := make([]models.Notification, len(userIDs))
	for i, userID := range userIDs {
		notifications[i] = models.Notification{
			UserID:  userID,
			Title:   title,
			Message: message,
			Type:    notificationType,
			IsRead:  false,
			Data:    models.NotificationData(data),
		}
	}

	return s.DB.Create(&notifications).Error
}
