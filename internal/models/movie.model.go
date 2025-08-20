package models

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Duration uint    `json:"duration"`
	MovieGenres []MovieGenre  `gorm:"foreignKey:MovieID" json:"movie_genres"`
	Schedules   []Schedule   `gorm:"foreignKey:MovieID" json:"schedules"`
}
