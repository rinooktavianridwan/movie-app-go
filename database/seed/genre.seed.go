package seed

import (
	"movie-app-go/internal/models"

	"gorm.io/gorm"
)

func SeedGenres(db *gorm.DB) ([]models.Genre, error) {
    genres := []models.Genre{
        {Name: "Action"},
        {Name: "Adventure"},
        {Name: "Animation"},
        {Name: "Comedy"},
        {Name: "Crime"},
        {Name: "Drama"},
        {Name: "Fantasy"},
        {Name: "Horror"},
        {Name: "Mystery"},
        {Name: "Romance"},
        {Name: "Sci-Fi"},
        {Name: "Thriller"},
        {Name: "War"},
        {Name: "Western"},
        {Name: "Biography"},
        {Name: "Documentary"},
        {Name: "Family"},
        {Name: "Musical"},
        {Name: "Sport"},
        {Name: "Superhero"},
    }
    
    if err := db.Create(&genres).Error; err != nil {
        return nil, err
    }
    return genres, nil
}