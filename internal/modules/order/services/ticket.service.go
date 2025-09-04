package services

import (
	"movie-app-go/internal/enums"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/order/repositories"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"

	"gorm.io/gorm"
)

type TicketService struct {
	TicketRepo *repositories.TicketRepository
}

func NewTicketService(ticketRepo *repositories.TicketRepository) *TicketService {
	return &TicketService{
		TicketRepo: ticketRepo,
	}
}

func (s *TicketService) GetTicketsByUser(userID uint, page, perPage int) (repository.PaginationResult[models.Ticket], error) {
	return s.TicketRepo.GetByUserIDPaginated(userID, page, perPage)
}

func (s *TicketService) GetAllTickets(page, perPage int) (repository.PaginationResult[models.Ticket], error) {
	return s.TicketRepo.GetAllPaginated(page, perPage)
}

func (s *TicketService) GetTicketByID(id uint, userID *uint) (*models.Ticket, error) {
	ticket, err := s.TicketRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrTicketNotFound
		}
		return nil, err
	}

	if userID != nil {
		ownershipCount, err := s.TicketRepo.CheckTicketOwnership(id, *userID)
		if err != nil {
			return nil, err
		}
		if ownershipCount == 0 {
			return nil, utils.ErrTicketNotFound
		}
	}

	return ticket, nil
}

func (s *TicketService) GetTicketsBySchedule(scheduleID uint, page, perPage int) (repository.PaginationResult[models.Ticket], error) {
	return s.TicketRepo.GetByScheduleIDPaginated(scheduleID, page, perPage)
}

func (s *TicketService) ScanTicket(id uint) error {
	var ticket *models.Ticket

	err := s.TicketRepo.WithTransaction(func(tx *gorm.DB) error {
        var err error
        ticket, err = s.TicketRepo.GetByID(id)
        if err != nil {
            if err == gorm.ErrRecordNotFound {
                return utils.ErrTicketNotFound
            }
            return err
        }

        switch ticket.Status {
        case enums.TicketStatusPending:
            return utils.ErrTicketNotPaid
        case enums.TicketStatusCancelled:
            return utils.ErrTicketCancelled
        case enums.TicketStatusUsed:
            return utils.ErrTicketAlreadyScanned
        case enums.TicketStatusActive:
            break
        default:
            return utils.ErrTicketNotFound
        }

        ticket.Status = enums.TicketStatusUsed
        return s.TicketRepo.UpdateTicket(ticket)
    })

    return err
}
