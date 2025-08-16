package seed

import (
    "movie-app-go/internal/models"
    "gorm.io/gorm"
)

func SeedFacilities(db *gorm.DB) ([]models.Facility, error) {
    facilities := []models.Facility{
        {Name: "Dolby Atmos"},
        {Name: "IMAX"},
        {Name: "4DX"},
        {Name: "VIP Seat"},
    }
    if err := db.Create(&facilities).Error; err != nil {
        return nil, err
    }
    return facilities, nil
}