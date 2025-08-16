package seed

import (
    "movie-app-go/internal/models"
    "gorm.io/gorm"
)

func SeedStudios(db *gorm.DB) ([]models.Studio, error) {
    studios := []models.Studio{
        {Name: "Studio 1", SeatCapacity: 100},
        {Name: "Studio 2", SeatCapacity: 80},
        {Name: "Studio 3", SeatCapacity: 120},
    }
    if err := db.Create(&studios).Error; err != nil {
        return nil, err
    }
    return studios, nil
}