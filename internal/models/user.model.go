package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string  `json:"name" gorm:"not null"`
	Email    string  `json:"email" gorm:"uniqueIndex;not null"`
	Password string  `json:"-" gorm:"not null"`
	Avatar   *string `json:"avatar" gorm:"default:null"`
	RoleID   *uint   `json:"role_id"`
	Role     *Role   `json:"role,omitempty" gorm:"foreignKey:RoleID"`

	Transactions []Transaction `gorm:"foreignKey:UserID" json:"transactions,omitempty"`
}
