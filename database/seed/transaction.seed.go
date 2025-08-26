package seed

import (
	"log"
	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"
	"time"

	"gorm.io/gorm"
)

func SeedTransactions(db *gorm.DB) ([]models.Transaction, error) {
	log.Println("Seeding transactions...")

	// Get users for relationships
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	if len(users) == 0 {
		log.Println("No users found, skipping transaction seeding")
		return []models.Transaction{}, nil
	}

	// Variasi user ID (cycle through available users)
	getUserID := func(index int) uint {
		return users[index%len(users)].ID
	}

	// Create transactions with specific timestamps
	transactionData := []struct {
		Transaction models.Transaction
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}{
		// August 2025 transactions - SUCCESS
		{
			Transaction: models.Transaction{
				UserID: getUserID(0), TotalAmount: 225000, PaymentMethod: "credit_card",
				PaymentStatus: constants.PaymentStatusSuccess,
			},
			CreatedAt: time.Date(2025, 8, 20, 15, 30, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 8, 20, 15, 35, 0, 0, time.UTC),
		},
		{
			Transaction: models.Transaction{
				UserID: getUserID(1), TotalAmount: 150000, PaymentMethod: "bank_transfer",
				PaymentStatus: constants.PaymentStatusSuccess,
			},
			CreatedAt: time.Date(2025, 8, 21, 10, 15, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 8, 21, 10, 20, 0, 0, time.UTC),
		},
		{
			Transaction: models.Transaction{
				UserID: getUserID(2), TotalAmount: 320000, PaymentMethod: "e_wallet",
				PaymentStatus: constants.PaymentStatusSuccess,
			},
			CreatedAt: time.Date(2025, 8, 21, 17, 45, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 8, 21, 17, 50, 0, 0, time.UTC),
		},
		{
			Transaction: models.Transaction{
				UserID: getUserID(0), TotalAmount: 160000, PaymentMethod: "credit_card",
				PaymentStatus: constants.PaymentStatusSuccess,
			},
			CreatedAt: time.Date(2025, 8, 22, 12, 20, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 8, 22, 12, 25, 0, 0, time.UTC),
		},

		// September 2025 transactions - SUCCESS
		{
			Transaction: models.Transaction{
				UserID: getUserID(1), TotalAmount: 225000, PaymentMethod: "bank_transfer",
				PaymentStatus: constants.PaymentStatusSuccess,
			},
			CreatedAt: time.Date(2025, 9, 1, 16, 10, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 9, 1, 16, 15, 0, 0, time.UTC),
		},
		{
			Transaction: models.Transaction{
				UserID: getUserID(2), TotalAmount: 240000, PaymentMethod: "e_wallet",
				PaymentStatus: constants.PaymentStatusSuccess,
			},
			CreatedAt: time.Date(2025, 9, 2, 19, 30, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 9, 2, 19, 35, 0, 0, time.UTC),
		},

		// October 2025 transactions - SUCCESS
		{
			Transaction: models.Transaction{
				UserID: getUserID(0), TotalAmount: 170000, PaymentMethod: "credit_card",
				PaymentStatus: constants.PaymentStatusSuccess,
			},
			CreatedAt: time.Date(2025, 10, 10, 20, 15, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 10, 10, 20, 20, 0, 0, time.UTC),
		},
		{
			Transaction: models.Transaction{
				UserID: getUserID(1), TotalAmount: 270000, PaymentMethod: "bank_transfer",
				PaymentStatus: constants.PaymentStatusSuccess,
			},
			CreatedAt: time.Date(2025, 10, 15, 22, 10, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 10, 15, 22, 15, 0, 0, time.UTC),
		},

		// Some FAILED transactions for completeness
		{
			Transaction: models.Transaction{
				UserID: getUserID(2), TotalAmount: 75000, PaymentMethod: "credit_card",
				PaymentStatus: constants.PaymentStatusFailed,
			},
			CreatedAt: time.Date(2025, 8, 23, 11, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 8, 23, 11, 5, 0, 0, time.UTC),
		},
		{
			Transaction: models.Transaction{
				UserID: getUserID(0), TotalAmount: 150000, PaymentMethod: "e_wallet",
				PaymentStatus: constants.PaymentStatusPending,
			},
			CreatedAt: time.Date(2025, 9, 5, 14, 30, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 9, 5, 14, 30, 0, 0, time.UTC),
		},
	}

	createdTransactions := []models.Transaction{}
	for _, data := range transactionData {
		transaction := data.Transaction

		if err := db.Create(&transaction).Error; err != nil {
			log.Printf("Error creating transaction: %v", err)
			continue
		}

		// Update timestamps after creation
		db.Model(&transaction).Updates(map[string]interface{}{
			"created_at": data.CreatedAt,
			"updated_at": data.UpdatedAt,
		})

		createdTransactions = append(createdTransactions, transaction)
	}

	log.Printf("Successfully seeded %d transactions", len(createdTransactions))
	return createdTransactions, nil
}
