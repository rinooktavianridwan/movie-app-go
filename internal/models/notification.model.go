package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type NotificationData map[string]interface{}

func (nd NotificationData) Value() (driver.Value, error) {
	return json.Marshal(nd)
}

func (nd *NotificationData) Scan(value interface{}) error {
	if value == nil {
		*nd = make(NotificationData)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, nd)
	case string:
		return json.Unmarshal([]byte(v), nd)
	}
	return nil
}

type Notification struct {
	ID        uint             `json:"id" gorm:"primaryKey"`
	UserID    uint             `json:"user_id" gorm:"not null"`
	Title     string           `json:"title" gorm:"not null;size:255"`
	Message   string           `json:"message" gorm:"not null;type:text"`
	Type      string           `json:"type" gorm:"not null;size:50"`
	IsRead    bool             `json:"is_read" gorm:"default:false"`
	Data      NotificationData `json:"data,omitempty" gorm:"type:jsonb"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `json:"-" gorm:"index"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
