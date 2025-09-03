package repositories

import (
    "movie-app-go/internal/enums"
    "movie-app-go/internal/models"
    "movie-app-go/internal/repository"
    "gorm.io/gorm"
)

type TransactionRepository struct {
    DB *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
    return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.Transaction], error) {
    query := r.DB.Preload("User").
        Preload("Tickets").
        Preload("Tickets.Schedule").
        Preload("Tickets.Schedule.Movie").
        Preload("Tickets.Schedule.Studio").
        Preload("Promo").
        Order("created_at DESC")

    return repository.Paginate[models.Transaction](query, page, perPage)
}

func (r *TransactionRepository) GetByUserIDPaginated(userID uint, page, perPage int) (repository.PaginationResult[models.Transaction], error) {
    query := r.DB.Preload("User").
        Preload("Tickets").
        Preload("Tickets.Schedule").
        Preload("Tickets.Schedule.Movie").
        Preload("Tickets.Schedule.Studio").
        Preload("Promo").
        Where("user_id = ?", userID).
        Order("created_at DESC")

    return repository.Paginate[models.Transaction](query, page, perPage)
}

func (r *TransactionRepository) GetByID(id uint) (*models.Transaction, error) {
    query := r.DB.Preload("User").
        Preload("Tickets").
        Preload("Tickets.Schedule").
        Preload("Tickets.Schedule.Movie").
        Preload("Tickets.Schedule.Studio").
        Preload("Promo")

    var transaction models.Transaction
    if err := query.First(&transaction, id).Error; err != nil {
        return nil, err
    }
    return &transaction, nil
}

func (r *TransactionRepository) GetByIDWithUserFilter(id uint, userID uint) (*models.Transaction, error) {
    query := r.DB.Preload("User").
        Preload("Tickets").
        Preload("Tickets.Schedule").
        Preload("Tickets.Schedule.Movie").
        Preload("Tickets.Schedule.Studio").
        Preload("Promo").
        Where("user_id = ?", userID)

    var transaction models.Transaction
    if err := query.First(&transaction, id).Error; err != nil {
        return nil, err
    }
    return &transaction, nil
}

func (r *TransactionRepository) GetScheduleWithMovieAndStudio(scheduleID uint) (*models.Schedule, error) {
    var schedule models.Schedule
    if err := r.DB.Preload("Movie").Preload("Studio").First(&schedule, scheduleID).Error; err != nil {
        return nil, err
    }
    return &schedule, nil
}

func (r *TransactionRepository) CountExistingTickets(scheduleID uint, seatNumbers []uint) (int64, error) {
    var count int64
    err := r.DB.Model(&models.Ticket{}).
        Where("schedule_id = ? AND seat_number IN ? AND status != ?", 
            scheduleID, seatNumbers, enums.TicketStatusCancelled).
        Count(&count).Error
    return count, err
}

func (r *TransactionRepository) CreateTransaction(transaction *models.Transaction) error {
    return r.DB.Create(transaction).Error
}

func (r *TransactionRepository) CreatePromoUsage(promoUsage *models.PromoUsage) error {
    return r.DB.Create(promoUsage).Error
}

func (r *TransactionRepository) IncrementPromoUsageCount(promoID uint) error {
    return r.DB.Model(&models.Promo{}).Where("id = ?", promoID).
        UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

func (r *TransactionRepository) CreateTickets(tickets []models.Ticket) error {
    return r.DB.Create(&tickets).Error
}

func (r *TransactionRepository) UpdateTransaction(transaction *models.Transaction) error {
    return r.DB.Save(transaction).Error
}

func (r *TransactionRepository) UpdateTicketsByTransactionID(transactionID uint, status string) error {
    return r.DB.Model(&models.Ticket{}).
        Where("transaction_id = ?", transactionID).
        Update("status", status).Error
}

func (r *TransactionRepository) GetScheduleByTransactionID(transactionID uint) (*models.Schedule, error) {
    var schedule models.Schedule
    err := r.DB.Preload("Movie").
        Joins("JOIN tickets ON tickets.schedule_id = schedules.id").
        Where("tickets.transaction_id = ?", transactionID).
        First(&schedule).Error
    return &schedule, err
}

func (r *TransactionRepository) WithTransaction(fn func(*gorm.DB) error) error {
    return r.DB.Transaction(fn)
}