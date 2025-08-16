package models

import "gorm.io/gorm"

type FacilityStudio struct {
	gorm.Model
	FacilityID uint     `json:"facility_id"`
	StudioID   uint     `json:"studio_id"`
	Facility   Facility `gorm:"foreignKey:FacilityID"`
	Studio     Studio   `gorm:"foreignKey:StudioID"`
}
