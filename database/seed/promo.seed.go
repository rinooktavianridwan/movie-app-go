package seed

import (
	"log"
	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"
	"time"

	"gorm.io/gorm"
)

func SeedPromos(db *gorm.DB) error {
	log.Println("Seeding promos...")

	var movies []models.Movie
	if err := db.Find(&movies).Error; err != nil {
		return err
	}

	promos := []models.Promo{
		{
			Name:          "Weekend Special",
			Code:          "WEEKEND20",
			Description:   "20% off all movie tickets during weekend",
			DiscountType:  constants.DiscountTypePercentage,
			DiscountValue: 20.0,
			MinTickets:    2,
			MaxDiscount:   &[]float64{50000}[0],
			UsageLimit:    &[]int{100}[0],
			IsActive:      true,
			ValidFrom:     time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			ValidUntil:    time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			Name:          "First Movie Special",
			Code:          "FIRST25K",
			Description:   "Rp 25.000 discount for first movie",
			DiscountType:  constants.DiscountTypeFixedAmount,
			DiscountValue: 25000,
			MinTickets:    1,
			IsActive:      true,
			ValidFrom:     time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
			ValidUntil:    time.Date(2025, 9, 30, 23, 59, 59, 0, time.UTC),
		},
		{
			Name:          "Student Discount",
			Code:          "STUDENT30",
			Description:   "30% off for students",
			DiscountType:  constants.DiscountTypePercentage,
			DiscountValue: 30.0,
			MinTickets:    1,
			MaxDiscount:   &[]float64{75000}[0],
			UsageLimit:    &[]int{50}[0],
			IsActive:      true,
			ValidFrom:     time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			ValidUntil:    time.Date(2025, 10, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			Name:          "Expired Promo",
			Code:          "EXPIRED10",
			Description:   "This promo has expired",
			DiscountType:  constants.DiscountTypePercentage,
			DiscountValue: 10.0,
			MinTickets:    1,
			IsActive:      true,
			ValidFrom:     time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
			ValidUntil:    time.Date(2025, 7, 31, 23, 59, 59, 0, time.UTC),
		},
	}

	createdPromos := []models.Promo{}
	for _, promo := range promos {
		var existing models.Promo
		if err := db.Where("code = ?", promo.Code).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&promo).Error; err != nil {
				log.Printf("Error creating promo %s: %v", promo.Code, err)
				continue
			}
			createdPromos = append(createdPromos, promo)
		} else {
			createdPromos = append(createdPromos, existing)
		}
	}
	if len(movies) > 0 {
		if len(createdPromos) >= 2 {
			promoMovie := models.PromoMovie{
				PromoID: createdPromos[1].ID,
				MovieID: movies[0].ID,
			}
			db.FirstOrCreate(&promoMovie, models.PromoMovie{
				PromoID: promoMovie.PromoID,
				MovieID: promoMovie.MovieID,
			})
		}
	}

	log.Printf("Successfully seeded %d promos", len(createdPromos))
	return nil
}
