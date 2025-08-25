package seed

import (
	"log"
	"movie-app-go/internal/models"
	"time"

	"gorm.io/gorm"
)

func SeedSchedules(db *gorm.DB) ([]models.Schedule, error) {
	log.Println("Seeding schedules...")

	// Get movies and studios for relationships
	var movies []models.Movie
	var studios []models.Studio

	if err := db.Find(&movies).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&studios).Error; err != nil {
		return nil, err
	}

	if len(movies) == 0 || len(studios) == 0 {
		log.Println("No movies or studios found, skipping schedule seeding")
		return []models.Schedule{}, nil
	}

	schedules := []models.Schedule{
		// August 2025 schedules
		{
			MovieID:   movies[0].ID,
			StudioID:  studios[0].ID,
			StartTime: time.Date(2025, 8, 20, 14, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 8, 20, 17, 0, 0, 0, time.UTC), // 3 hours duration
			Date:      time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
			Price:     75000,
		},
		{
			MovieID:   movies[0].ID,
			StudioID:  studios[0].ID,
			StartTime: time.Date(2025, 8, 20, 18, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 8, 20, 21, 0, 0, 0, time.UTC),
			Date:      time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
			Price:     75000,
		},
		{
			MovieID:   movies[0].ID,
			StudioID:  studios[1].ID,
			StartTime: time.Date(2025, 8, 21, 16, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 8, 21, 19, 0, 0, 0, time.UTC),
			Date:      time.Date(2025, 8, 21, 0, 0, 0, 0, time.UTC),
			Price:     75000,
		},
		{
			MovieID:   movies[0].ID,
			StudioID:  studios[1].ID,
			StartTime: time.Date(2025, 8, 22, 20, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 8, 22, 23, 0, 0, 0, time.UTC),
			Date:      time.Date(2025, 8, 22, 0, 0, 0, 0, time.UTC),
			Price:     75000,
		},

		// Different movies - August (use last movie in array)
		{
			MovieID:   movies[len(movies)-1].ID,
			StudioID:  studios[0].ID,
			StartTime: time.Date(2025, 8, 21, 12, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 8, 21, 15, 0, 0, 0, time.UTC),
			Date:      time.Date(2025, 8, 21, 0, 0, 0, 0, time.UTC),
			Price:     80000,
		},
		{
			MovieID:   movies[len(movies)-1].ID,
			StudioID:  studios[1].ID,
			StartTime: time.Date(2025, 8, 22, 15, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 8, 22, 18, 0, 0, 0, time.UTC),
			Date:      time.Date(2025, 8, 22, 0, 0, 0, 0, time.UTC),
			Price:     80000,
		},

		// September 2025 schedules
		{
			MovieID:   movies[0].ID,
			StudioID:  studios[0].ID,
			StartTime: time.Date(2025, 9, 1, 14, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 9, 1, 17, 0, 0, 0, time.UTC),
			Date:      time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
			Price:     75000,
		},
		{
			MovieID:   movies[0].ID,
			StudioID:  studios[1].ID,
			StartTime: time.Date(2025, 9, 2, 18, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 9, 2, 21, 0, 0, 0, time.UTC),
			Date:      time.Date(2025, 9, 2, 0, 0, 0, 0, time.UTC),
			Price:     75000,
		},
		{
			MovieID:   movies[len(movies)-1].ID,
			StudioID:  studios[0].ID,
			StartTime: time.Date(2025, 9, 3, 16, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 9, 3, 19, 0, 0, 0, time.UTC),
			Date:      time.Date(2025, 9, 3, 0, 0, 0, 0, time.UTC),
			Price:     80000,
		},

		// October 2025 schedules
		{
			MovieID:   movies[0].ID,
			StudioID:  studios[0].ID,
			StartTime: time.Date(2025, 10, 10, 19, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 10, 10, 22, 0, 0, 0, time.UTC),
			Date:      time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC),
			Price:     85000,
		},
		{
			MovieID:   movies[len(movies)-1].ID,
			StudioID:  studios[1].ID,
			StartTime: time.Date(2025, 10, 15, 21, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2025, 10, 16, 0, 0, 0, 0, time.UTC),
			Date:      time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
			Price:     90000,
		},
	}

	// Check existing schedules to avoid duplicates
	for _, schedule := range schedules {
		var existing models.Schedule
		err := db.Where("movie_id = ? AND studio_id = ? AND start_time = ?",
			schedule.MovieID, schedule.StudioID, schedule.StartTime).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&schedule).Error; err != nil {
				return nil, err
			}
		}
	}

	// Return created schedules
	var createdSchedules []models.Schedule
	if err := db.Find(&createdSchedules).Error; err != nil {
		return nil, err
	}

	log.Printf("Successfully seeded %d schedules", len(createdSchedules))
	return createdSchedules, nil
}
