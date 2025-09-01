package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Avatar   *string `json:"avatar" gorm:"default:null"`
	IsAdmin  bool   `json:"is_admin"`

	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
}
