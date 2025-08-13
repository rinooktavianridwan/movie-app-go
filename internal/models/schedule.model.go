package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	MovieID   uint      `json:"movie_id"`
	StudioID  uint      `json:"studio_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Date      time.Time `json:"date"`
	Price     float64   `json:"price"`
}
