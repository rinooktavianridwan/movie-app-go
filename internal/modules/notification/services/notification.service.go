package services

import (
	"fmt"
	"log"
	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/notification/options"
	"movie-app-go/internal/modules/notification/repositories"
	"movie-app-go/internal/modules/notification/requests"
	"movie-app-go/internal/modules/notification/responses"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"

	"gorm.io/gorm"
)

type NotificationService struct {
	NotificationRepo *repositories.NotificationRepository
}

func NewNotificationService(notificationRepo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{
		NotificationRepo: notificationRepo,
	}
}

func (s *NotificationService) CreateNotification(req *requests.CreateNotificationRequest) error {
	notification := models.Notification{
		UserID:  req.UserID,
		Title:   req.Title,
		Message: req.Message,
		Type:    req.Type,
		IsRead:  false,
		Data:    req.Data,
	}

	return s.NotificationRepo.Create(&notification)
}

func (s *NotificationService) GetUserNotifications(userID uint, opts options.NotificationOptions) (repository.PaginationResult[models.Notification], error) {
	return s.NotificationRepo.GetByUserIDWithOptions(userID, opts)
}

func (s *NotificationService) GetNotificationByID(id uint) (*models.Notification, error) {
	notification, err := s.NotificationRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotificationNotFound
		}
		return nil, err
	}
	return notification, nil
}

func (s *NotificationService) MarkAsRead(id uint, userID uint) error {
	notification, err := s.NotificationRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrNotificationNotFound
		}
		return err
	}

	if notification.UserID != userID {
		return utils.ErrNotificationNotFound
	}

	notification.IsRead = true
	return s.NotificationRepo.Update(notification)
}

func (s *NotificationService) MarkAsUnread(id uint, userID uint) error {
	notification, err := s.NotificationRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrNotificationNotFound
		}
		return err
	}

	if notification.UserID != userID {
		return utils.ErrNotificationNotFound
	}

	notification.IsRead = false
	return s.NotificationRepo.Update(notification)
}

func (s *NotificationService) BulkMarkAsRead(req *requests.BulkMarkReadRequest, userID uint) error {
	if len(req.NotificationIDs) == 0 {
		return fmt.Errorf("no notification IDs provided")
	}

	rowsAffected, err := s.NotificationRepo.UpdateByUserIDAndIDs(
		userID,
		req.NotificationIDs,
		map[string]interface{}{"is_read": true},
	)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return utils.ErrNotificationNotFound
	}

	return nil
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.NotificationRepo.UpdateByUserID(userID, map[string]interface{}{"is_read": true})
}

func (s *NotificationService) GetNotificationStats(userID uint) (*responses.NotificationStatsResponse, error) {
	totalCount, err := s.NotificationRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}

	unreadCount, err := s.NotificationRepo.CountUnreadByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &responses.NotificationStatsResponse{
		TotalCount:  totalCount,
		UnreadCount: unreadCount,
		ReadCount:   totalCount - unreadCount,
	}, nil
}

func (s *NotificationService) DeleteNotification(id uint, userID uint) error {
	rowsAffected, err := s.NotificationRepo.Delete(id, userID)
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return utils.ErrNotificationNotFound
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

    return s.CreateNotification(req)
}

func (s *NotificationService) CreatePromoNotification(promoID uint, promoName, promoCode string) error {
	userIDs, err := s.NotificationRepo.GetAllUserIDs()
    if err != nil {
        log.Printf("Failed to get users for promo notification: %v", err)
        return err
    }

    if len(userIDs) == 0 {
        return utils.ErrUserNotFound
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

    return s.CreateNotification(req)
}

func (s *NotificationService) CreateBulkNotifications(
	userIDs []uint,
	title string,
	message string,
	notificationType string,
	data map[string]interface{},
) error {
	if len(userIDs) == 0 {
        return utils.ErrUserNotFound
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

    return s.NotificationRepo.CreateBatch(notifications)
}

func (s *NotificationService) CreateSystemNotification(req *requests.CreateSystemNotificationRequest) error {
    userIDs, err := s.NotificationRepo.GetAllUserIDs()
    if err != nil {
        return fmt.Errorf("failed to get users: %w", err)
    }

    if len(userIDs) == 0 {
        return utils.ErrUserNotFound
    }

    return s.CreateBulkNotifications(
        userIDs,
        req.Title,
        req.Message,
        constants.NotificationTypeSystem,
        req.Data,
    )
}
