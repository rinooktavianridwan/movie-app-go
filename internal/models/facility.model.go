package models

import "gorm.io/gorm"

type Facility struct {
	gorm.Model
	Name            string           `json:"name"`
	FacilityStudios []FacilityStudio `gorm:"foreignKey:FacilityID"`
}
