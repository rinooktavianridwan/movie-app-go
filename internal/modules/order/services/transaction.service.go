package services

import (
	"fmt"
	"movie-app-go/internal/constants"
	"movie-app-go/internal/jobs"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/order/requests"
	"movie-app-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type TransactionService struct {
	DB           *gorm.DB
	QueueService *jobs.QueueService
}

func NewTransactionService(db *gorm.DB, queueService *jobs.QueueService) *TransactionService {
	return &TransactionService{DB: db, QueueService: queueService}
}

func (s *TransactionService) CreateTransaction(userID uint, req *requests.CreateTransactionRequest) (*models.Transaction, error) {
	var transaction *models.Transaction

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var schedule models.Schedule
		if err := tx.Preload("Movie").Preload("Studio").First(&schedule, req.ScheduleID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("schedule not found")
			}
			return err
		}

		if len(req.SeatNumbers) > int(schedule.Studio.SeatCapacity) {
			return fmt.Errorf("seat numbers exceed studio capacity")
		}

		var existingTickets int64
		if err := tx.Model(&models.Ticket{}).
			Where("schedule_id = ? AND seat_number IN ? AND status != ?",
				req.ScheduleID, req.SeatNumbers, constants.TicketStatusCancelled).
			Count(&existingTickets).Error; err != nil {
			return err
		}
		if existingTickets > 0 {
			return fmt.Errorf("some seats are already booked")
		}

		for _, seatNum := range req.SeatNumbers {
			if seatNum < 1 || seatNum > schedule.Studio.SeatCapacity {
				return fmt.Errorf("invalid seat number: %d", seatNum)
			}
		}

		totalAmount := float64(len(req.SeatNumbers)) * schedule.Price

		transaction = &models.Transaction{
			UserID:        userID,
			TotalAmount:   totalAmount,
			PaymentMethod: req.PaymentMethod,
			PaymentStatus: constants.PaymentStatusPending,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		tickets := make([]models.Ticket, 0, len(req.SeatNumbers))
		for _, seatNum := range req.SeatNumbers {
			tickets = append(tickets, models.Ticket{
				TransactionID: transaction.ID,
				ScheduleID:    req.ScheduleID,
				SeatNumber:    seatNum,
				Status:        constants.TicketStatusPending,
				Price:         schedule.Price,
			})
		}

		if err := tx.Create(&tickets).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := s.QueueService.SchedulePaymentTimeout(transaction.ID, 2*time.Minute); err != nil {
		fmt.Printf("Failed to schedule payment timeout job: %v\n", err)
	}

	return transaction, nil
}

func (s *TransactionService) GetTransactionsByUser(userID uint, page, perPage int) (repository.PaginationResult[models.Transaction], error) {
	query := s.DB.Preload("User").
		Preload("Tickets").
		Preload("Tickets.Schedule").
		Preload("Tickets.Schedule.Movie").
		Preload("Tickets.Schedule.Studio").
		Where("user_id = ?", userID).
		Order("created_at DESC")

	return repository.Paginate[models.Transaction](query, page, perPage)
}

func (s *TransactionService) GetAllTransactions(page, perPage int) (repository.PaginationResult[models.Transaction], error) {
	query := s.DB.Preload("User").
		Preload("Tickets").
		Preload("Tickets.Schedule").
		Preload("Tickets.Schedule.Movie").
		Preload("Tickets.Schedule.Studio").
		Order("created_at DESC")

	return repository.Paginate[models.Transaction](query, page, perPage)
}

func (s *TransactionService) GetTransactionByID(id uint, userID *uint) (*models.Transaction, error) {
	query := s.DB.Preload("User").
		Preload("Tickets").
		Preload("Tickets.Schedule").
		Preload("Tickets.Schedule.Movie").
		Preload("Tickets.Schedule.Studio")

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	var transaction models.Transaction
	if err := query.First(&transaction, id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (s *TransactionService) ProcessPayment(id uint, req *requests.ProcessPaymentRequest) (*models.Transaction, error) {
	var transaction models.Transaction

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&transaction, id).Error; err != nil {
			return err
		}

		if transaction.PaymentStatus != constants.PaymentStatusPending {
			return fmt.Errorf("transaction status can only be updated from pending")
		}

		transaction.PaymentStatus = req.PaymentStatus
		if err := tx.Save(&transaction).Error; err != nil {
			return err
		}

		if req.PaymentStatus == constants.PaymentStatusSuccess {
            if err := tx.Model(&models.Ticket{}).
                Where("transaction_id = ?", id).
                Update("status", constants.TicketStatusActive).Error; err != nil {
                return err
            }
        } else if req.PaymentStatus == constants.PaymentStatusFailed {
            if err := tx.Model(&models.Ticket{}).
                Where("transaction_id = ?", id).
                Update("status", constants.TicketStatusCancelled).Error; err != nil {
                return err
            }
        }

		return nil
	})

	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
