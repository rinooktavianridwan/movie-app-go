package services

import (
	"fmt"
	"log"
	"movie-app-go/internal/constants"
	"movie-app-go/internal/jobs"
	"movie-app-go/internal/models"
	notificationServices "movie-app-go/internal/modules/notification/services"
	"movie-app-go/internal/modules/order/requests"
	promoRequests "movie-app-go/internal/modules/promo/requests"
	promoServices "movie-app-go/internal/modules/promo/services"
	"movie-app-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type TransactionService struct {
	DB                  *gorm.DB
	QueueService        *jobs.QueueService
	PromoService        *promoServices.PromoService
	NotificationService *notificationServices.NotificationService
}

func NewTransactionService(db *gorm.DB, queueService *jobs.QueueService, promoService *promoServices.PromoService) *TransactionService {
	return &TransactionService{DB: db, QueueService: queueService, PromoService: promoService, NotificationService: notificationServices.NewNotificationService(db)}
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

		originalAmount := float64(len(req.SeatNumbers)) * schedule.Price
		finalAmount := originalAmount
		discountAmount := float64(0)
		var promoID *uint

		if req.PromoCode != "" && s.PromoService != nil {
			promoValidation := &promoRequests.ValidatePromoRequest{
				PromoCode:   req.PromoCode,
				TotalAmount: originalAmount,
				MovieIDs:    []uint{schedule.MovieID},
				SeatNumbers: req.SeatNumbers,
			}

			result, err := s.PromoService.ValidatePromo(userID, promoValidation)
			if err != nil || !result.IsValid {
				return fmt.Errorf("invalid promo: %s", result.Message)
			}

			discountAmount = result.DiscountAmount
			finalAmount = result.FinalAmount

			promo, err := s.PromoService.GetPromoByCode(req.PromoCode)
			if err == nil {
				promoID = &promo.ID
			}
		}

		transaction = &models.Transaction{
			UserID:         userID,
			TotalAmount:    finalAmount,
			OriginalAmount: &originalAmount,
			DiscountAmount: discountAmount,
			PaymentMethod:  req.PaymentMethod,
			PaymentStatus:  constants.PaymentStatusPending,
			PromoID:        promoID,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		if promoID != nil {
			promoUsage := models.PromoUsage{
				PromoID:        *promoID,
				UserID:         userID,
				TransactionID:  transaction.ID,
				DiscountAmount: discountAmount,
			}
			if err := tx.Create(&promoUsage).Error; err != nil {
				return err
			}

			if err := tx.Model(&models.Promo{}).Where("id = ?", *promoID).
				UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1)).Error; err != nil {
				return err
			}
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
		Preload("Promo").
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
		Preload("Promo").
		Order("created_at DESC")

	return repository.Paginate[models.Transaction](query, page, perPage)
}

func (s *TransactionService) GetTransactionByID(id uint, userID *uint) (*models.Transaction, error) {
	query := s.DB.Preload("User").
		Preload("Tickets").
		Preload("Tickets.Schedule").
		Preload("Tickets.Schedule.Movie").
		Preload("Tickets.Schedule.Studio").
		Preload("Promo")

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
			return fmt.Errorf("transaction already failed or timed out")
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

	if req.PaymentStatus == constants.PaymentStatusSuccess {
		var schedule models.Schedule
		if err := s.DB.Preload("Movie").
			Joins("JOIN tickets ON tickets.schedule_id = schedules.id").
			Where("tickets.transaction_id = ?", id).
			First(&schedule).Error; err != nil {
			log.Printf("⚠️ Failed to get schedule for transaction %d: %v", id, err)
		} else {
			// ✅ CLEAN: Use notification service directly
			go func() {
				if err := s.NotificationService.CreateBookingConfirmationNotification(
					transaction.UserID,
					transaction.ID,
					schedule.Movie.Title,
					transaction.TotalAmount,
				); err != nil {
					log.Printf("Failed to create booking confirmation: %v", err)
				}
			}()

			go func() {
				if err := s.NotificationService.CreateMovieReminderNotification(
					transaction.UserID,
					schedule.Movie.Title,
					schedule.StartTime.Format("2006-01-02 15:04:05"),
					schedule.MovieID,
				); err != nil {
					log.Printf("Failed to create movie reminder: %v", err)
				}
			}()
		}
	}

	return &transaction, nil
}
