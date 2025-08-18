package models

import "gorm.io/gorm"

type Ticket struct {
	gorm.Model
	TransactionID uint    `gorm:"not null" json:"transaction_id"`
	ScheduleID    uint    `gorm:"not null" json:"schedule_id"`
	SeatNumber    uint    `json:"seat_number"`
	Status        string  `json:"status"`
	Price         float64 `gorm:"not null" json:"price"`

	Transaction Transaction `gorm:"foreignKey:TransactionID" json:"transaction"`
	Schedule    Schedule    `gorm:"foreignKey:ScheduleID" json:"schedule"`
}
