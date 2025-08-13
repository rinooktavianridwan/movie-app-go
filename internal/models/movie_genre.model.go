package models

import "gorm.io/gorm"

type MovieGenre struct {
	gorm.Model
	MovieID uint `json:"movie_id"`
	GenreID uint `json:"genre_id"`
}
