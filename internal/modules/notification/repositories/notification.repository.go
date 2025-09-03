package repositories

import (
    "movie-app-go/internal/models"
    "movie-app-go/internal/modules/notification/options"
    "movie-app-go/internal/repository"
    "gorm.io/gorm"
)

type NotificationRepository struct {
    DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
    return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) GetByUserIDWithOptions(userID uint, opts options.NotificationOptions) (repository.PaginationResult[models.Notification], error) {
    query := r.DB.Model(&models.Notification{}).Where("user_id = ?", userID)

    if opts.Type != "" {
        query = query.Where("type = ?", opts.Type)
    }

    if opts.IsRead != nil {
        query = query.Where("is_read = ?", *opts.IsRead)
    }

    query = query.Order("created_at DESC")

    return repository.Paginate[models.Notification](query, opts.Page, opts.PerPage)
}

func (r *NotificationRepository) GetByID(id uint) (*models.Notification, error) {
    var notification models.Notification
    if err := r.DB.First(&notification, id).Error; err != nil {
        return nil, err
    }
    return &notification, nil
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
    return r.DB.Create(notification).Error
}

func (r *NotificationRepository) CreateBatch(notifications []models.Notification) error {
    if len(notifications) == 0 {
        return nil
    }
    return r.DB.Create(&notifications).Error
}

func (r *NotificationRepository) Update(notification *models.Notification) error {
    return r.DB.Save(notification).Error
}

func (r *NotificationRepository) Delete(id, userID uint) (int64, error) {
    result := r.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{})
    return result.RowsAffected, result.Error
}

func (r *NotificationRepository) UpdateByUserIDAndIDs(userID uint, notificationIDs []uint, updates map[string]interface{}) (int64, error) {
    result := r.DB.Model(&models.Notification{}).
        Where("id IN ? AND user_id = ?", notificationIDs, userID).
        Updates(updates)
    return result.RowsAffected, result.Error
}

func (r *NotificationRepository) UpdateByUserID(userID uint, updates map[string]interface{}) error {
    return r.DB.Model(&models.Notification{}).
        Where("user_id = ? AND is_read = ?", userID, false).
        Updates(updates).Error
}

func (r *NotificationRepository) CountByUserID(userID uint) (int64, error) {
    var count int64
    err := r.DB.Model(&models.Notification{}).Where("user_id = ?", userID).Count(&count).Error
    return count, err
}

func (r *NotificationRepository) CountUnreadByUserID(userID uint) (int64, error) {
    var count int64
    err := r.DB.Model(&models.Notification{}).
        Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
    return count, err
}

func (r *NotificationRepository) GetAllUserIDs() ([]uint, error) {
    var userIDs []uint
    err := r.DB.Model(&models.User{}).
        Where("is_admin = ?", false).
        Pluck("id", &userIDs).Error
    return userIDs, err
}
