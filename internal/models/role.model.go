package models

import "gorm.io/gorm"

type Role struct {
    gorm.Model
    Name        string              `json:"name" gorm:"uniqueIndex;not null"`
    Description string              `json:"description"`
	
    Permissions []Permission        `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
    Users       []User              `json:"users,omitempty" gorm:"foreignKey:RoleID"`
}