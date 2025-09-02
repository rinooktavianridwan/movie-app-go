package services

import (
	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"

	"gorm.io/gorm"
)

type TicketService struct {
	DB *gorm.DB
}

func NewTicketService(db *gorm.DB) *TicketService {
	return &TicketService{DB: db}
}

func (s *TicketService) GetTicketsByUser(userID uint, page, perPage int) (repository.PaginationResult[models.Ticket], error) {
	query := s.DB.Preload("Transaction").
		Preload("Transaction.User").
		Preload("Schedule").
		Preload("Schedule.Movie").
		Preload("Schedule.Studio").
		Joins("JOIN transactions ON tickets.transaction_id = transactions.id").
		Where("transactions.user_id = ?", userID).
		Order("tickets.created_at DESC")

	return repository.Paginate[models.Ticket](query, page, perPage)
}

func (s *TicketService) GetAllTickets(page, perPage int) (repository.PaginationResult[models.Ticket], error) {
	query := s.DB.Preload("Transaction").
		Preload("Transaction.User").
		Preload("Schedule").
		Preload("Schedule.Movie").
		Preload("Schedule.Studio").
		Order("created_at DESC")

	return repository.Paginate[models.Ticket](query, page, perPage)
}

func (s *TicketService) GetTicketByID(id uint, userID *uint) (*models.Ticket, error) {
	query := s.DB.Preload("Transaction").
		Preload("Transaction.User").
		Preload("Schedule").
		Preload("Schedule.Movie").
		Preload("Schedule.Studio")

	if userID != nil {
		query = query.Joins("JOIN transactions ON tickets.transaction_id = transactions.id").
			Where("transactions.user_id = ?", *userID)
	}

	var ticket models.Ticket
	if err := query.First(&ticket, id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (s *TicketService) GetTicketsBySchedule(scheduleID uint, page, perPage int) (repository.PaginationResult[models.Ticket], error) {
	query := s.DB.Preload("Transaction").
		Preload("Transaction.User").
		Preload("Schedule").
		Preload("Schedule.Movie").
		Preload("Schedule.Studio").
		Where("schedule_id = ?", scheduleID).
		Order("seat_number ASC")

	return repository.Paginate[models.Ticket](query, page, perPage)
}

func (s *TicketService) ScanTicket(id uint, userID *uint) error {
	var ticket models.Ticket

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		query := tx.Preload("Transaction").
			Preload("Schedule").
			Preload("Schedule.Movie").
			Preload("Schedule.Studio")

		if userID != nil {
			query = query.Joins("JOIN transactions ON tickets.transaction_id = transactions.id").
				Where("transactions.user_id = ?", *userID)
		}

		if err := query.First(&ticket, id).Error; err != nil {
			return err
		}

		switch ticket.Status {
		case constants.TicketStatusPending:
			return utils.ErrTicketNotPaid
		case constants.TicketStatusCancelled:
			return utils.ErrTicketCancelled
		case constants.TicketStatusUsed:
			return utils.ErrTicketAlreadyScanned
		case constants.TicketStatusActive:
			break
		default:
			return utils.ErrTicketNotFound
		}

		ticket.Status = constants.TicketStatusUsed
		if err := tx.Save(&ticket).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
