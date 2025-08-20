package models

import "gorm.io/gorm"

type Studio struct {
	gorm.Model
	Name            string           `json:"name"`
	SeatCapacity    uint             `json:"seat_capacity"`
	FacilityStudios []FacilityStudio `gorm:"foreignKey:StudioID"`
	Schedules       []Schedule       `gorm:"foreignKey:StudioID" json:"schedules"`
}
