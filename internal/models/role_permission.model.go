package models

import "gorm.io/gorm"

type RolePermission struct {
	gorm.Model
	RoleID       uint       `json:"role_id" gorm:"not null"`
	PermissionID uint       `json:"permission_id" gorm:"not null"`
	Role         Role       `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Permission   Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID"`
}
