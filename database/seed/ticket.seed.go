package seed

import (
	"log"
	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"

	"gorm.io/gorm"
)

func SeedTickets(db *gorm.DB, transactions []models.Transaction, schedules []models.Schedule) error {
	log.Println("Seeding tickets...")

	if len(transactions) == 0 || len(schedules) == 0 {
		log.Println("No transactions or schedules found, skipping ticket seeding")
		return nil
	}

	// Map untuk ticket seeding berdasarkan transaction dan schedule
	ticketConfigs := []struct {
		TransactionIndex int
		ScheduleIndex    int
		SeatNumbers      []int
		Price            float64
		Status           string
	}{
		// Transaction 1: 3 tickets (August 20)
		{0, 0, []int{5, 6, 7}, 75000, constants.TicketStatusUsed},

		// Transaction 2: 2 tickets (August 21)
		{1, 2, []int{1, 2}, 75000, constants.TicketStatusActive},

		// Transaction 3: 4 tickets (August 21)
		{2, 4, []int{10, 11, 12, 13}, 80000, constants.TicketStatusActive},

		// Transaction 4: 2 tickets (August 22)
		{3, 5, []int{8, 9}, 80000, constants.TicketStatusUsed},

		// Transaction 5: 3 tickets (September 1)
		{4, 6, []int{3, 4, 5}, 75000, constants.TicketStatusActive},

		// Transaction 6: 3 tickets (September 2)
		{5, 7, []int{15, 16, 17}, 75000, constants.TicketStatusActive},

		// Transaction 7: 2 tickets (October 10)
		{6, 9, []int{20, 21}, 85000, constants.TicketStatusActive},

		// Transaction 8: 3 tickets (October 15)
		{7, 10, []int{25, 26, 27}, 90000, constants.TicketStatusActive},
	}

	for _, config := range ticketConfigs {
		if config.TransactionIndex >= len(transactions) || config.ScheduleIndex >= len(schedules) {
			continue // Skip jika index tidak valid
		}

		transaction := transactions[config.TransactionIndex]
		schedule := schedules[config.ScheduleIndex]

		// Only create tickets for successful transactions
		if transaction.PaymentStatus != constants.PaymentStatusSuccess {
			continue
		}

		for _, seatNum := range config.SeatNumbers {
			ticket := models.Ticket{
				TransactionID: transaction.ID,
				ScheduleID:    schedule.ID,
				SeatNumber:    uint(seatNum),
				Status:        config.Status,
				Price:         config.Price,
			}

			// Check if ticket already exists
			var existing models.Ticket
			err := db.Where("transaction_id = ? AND schedule_id = ? AND seat_number = ?",
				ticket.TransactionID, ticket.ScheduleID, ticket.SeatNumber).First(&existing).Error

			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&ticket).Error; err != nil {
					log.Printf("Error creating ticket: %v", err)
					continue
				}
			}
		}
	}

	// Count created tickets
	var ticketCount int64
	db.Model(&models.Ticket{}).Count(&ticketCount)

	log.Printf("Successfully seeded tickets. Total tickets in DB: %d", ticketCount)
	return nil
}
