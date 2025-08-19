package models

import "gorm.io/gorm"

type MovieGenre struct {
	gorm.Model
	MovieID uint `json:"movie_id"`
	GenreID uint `json:"genre_id"`
	Movie   Movie  `gorm:"foreignKey:MovieID"`
    Genre   Genre  `gorm:"foreignKey:GenreID"`
}
