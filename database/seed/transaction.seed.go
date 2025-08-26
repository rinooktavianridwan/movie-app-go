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

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}

	if len(users) == 0 {
		log.Println("No users found, skipping transaction seeding")
		return []models.Transaction{}, nil
	}

	var promos []models.Promo
	db.Find(&promos)

	getUserID := func(index int) uint {
		return users[index%len(users)].ID
	}

	getPromoID := func(code string) *uint {
		for _, promo := range promos {
			if promo.Code == code {
				return &promo.ID
			}
		}
		return nil
	}

	transactionData := []struct {
		Transaction models.Transaction
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}{
		{
			Transaction: models.Transaction{
				UserID:         getUserID(0),
				OriginalAmount: &[]float64{225000}[0],
				DiscountAmount: 0,
				TotalAmount:    225000,
				PaymentMethod:  "credit_card",
				PaymentStatus:  constants.PaymentStatusSuccess,
				PromoID:        nil,
			},
			CreatedAt: time.Date(2025, 8, 20, 15, 30, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 8, 20, 15, 35, 0, 0, time.UTC),
		},
		{
			Transaction: models.Transaction{
				UserID:         getUserID(1),
				OriginalAmount: &[]float64{150000}[0],
				DiscountAmount: 0,
				TotalAmount:    150000,
				PaymentMethod:  "bank_transfer",
				PaymentStatus:  constants.PaymentStatusSuccess,
				PromoID:        nil,
			},
			CreatedAt: time.Date(2025, 8, 21, 10, 15, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 8, 21, 10, 20, 0, 0, time.UTC),
		},

		{
			Transaction: models.Transaction{
				UserID:         getUserID(0),
				OriginalAmount: &[]float64{225000}[0],
				DiscountAmount: 45000,
				TotalAmount:    180000,
				PaymentMethod:  "credit_card",
				PaymentStatus:  constants.PaymentStatusSuccess,
				PromoID:        getPromoID("WEEKEND20"),
			},
			CreatedAt: time.Date(2025, 8, 25, 14, 30, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 8, 25, 14, 35, 0, 0, time.UTC),
		},
		{
			Transaction: models.Transaction{
				UserID:         getUserID(1),
				OriginalAmount: &[]float64{150000}[0],
				DiscountAmount: 25000,
				TotalAmount:    125000,
				PaymentMethod:  "bank_transfer",
				PaymentStatus:  constants.PaymentStatusSuccess,
				PromoID:        getPromoID("FIRST25K"),
			},
			CreatedAt: time.Date(2025, 9, 5, 16, 20, 0, 0, time.UTC),
			UpdatedAt: time.Date(2025, 9, 5, 16, 25, 0, 0, time.UTC),
		},

	}

	createdTransactions := []models.Transaction{}
	for _, data := range transactionData {
		transaction := data.Transaction

		if err := db.Create(&transaction).Error; err != nil {
			log.Printf("Error creating transaction: %v", err)
			continue
		}

		db.Model(&transaction).Updates(map[string]interface{}{
			"created_at": data.CreatedAt,
			"updated_at": data.UpdatedAt,
		})

		if transaction.PromoID != nil && transaction.DiscountAmount > 0 {
			promoUsage := models.PromoUsage{
				PromoID:        *transaction.PromoID,
				UserID:         transaction.UserID,
				TransactionID:  transaction.ID,
				DiscountAmount: transaction.DiscountAmount,
				UsedAt:         data.CreatedAt,
			}
			db.Create(&promoUsage)


			db.Model(&models.Promo{}).Where("id = ?", *transaction.PromoID).
				UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1))
		}

		createdTransactions = append(createdTransactions, transaction)
	}

	log.Printf("Successfully seeded %d transactions", len(createdTransactions))
	return createdTransactions, nil
}
