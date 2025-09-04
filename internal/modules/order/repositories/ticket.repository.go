package repositories

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/repository"

	"gorm.io/gorm"
)

type TicketRepository struct {
	DB *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{DB: db}
}

func (r *TicketRepository) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.Ticket], error) {
	query := r.DB.Preload("Transaction").
		Preload("Transaction.User").
		Preload("Schedule").
		Preload("Schedule.Movie").
		Preload("Schedule.Studio")

	return repository.Paginate[models.Ticket](query, page, perPage)
}

func (r *TicketRepository) GetByUserIDPaginated(userID uint, page, perPage int) (repository.PaginationResult[models.Ticket], error) {
	query := r.DB.Joins("JOIN transactions ON tickets.transaction_id = transactions.id").
		Where("transactions.user_id = ?", userID).
		Preload("Transaction").
		Preload("Transaction.User").
		Preload("Schedule").
		Preload("Schedule.Movie").
		Preload("Schedule.Studio").
		Order("tickets.created_at DESC")

	return repository.Paginate[models.Ticket](query, page, perPage)
}

func (r *TicketRepository) GetByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := r.DB.Preload("Transaction").
		Preload("Transaction.User").
		Preload("Schedule").
		Preload("Schedule.Movie").
		Preload("Schedule.Studio").
		First(&ticket, id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) GetByScheduleIDPaginated(scheduleID uint, page, perPage int) (repository.PaginationResult[models.Ticket], error) {
	query := r.DB.Where("schedule_id = ?", scheduleID).
		Preload("Transaction").
		Preload("Transaction.User").
		Order("seat_number ASC")

	return repository.Paginate[models.Ticket](query, page, perPage)
}

func (r *TicketRepository) GetByScheduleID(scheduleID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.DB.Where("schedule_id = ?", scheduleID).
		Preload("Transaction").
		Preload("Transaction.User").
		Order("seat_number ASC").
		Find(&tickets).Error
	return tickets, err
}

func (r *TicketRepository) CheckTicketOwnership(ticketID, userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Ticket{}).
		Joins("JOIN transactions ON tickets.transaction_id = transactions.id").
		Where("tickets.id = ? AND transactions.user_id = ?", ticketID, userID).
		Count(&count).Error
	return count, err
}

func (r *TicketRepository) UpdateTicket(ticket *models.Ticket) error {
	return r.DB.Save(ticket).Error
}

func (r *TicketRepository) WithTransaction(fn func(*gorm.DB) error) error {
	return r.DB.Transaction(fn)
}
