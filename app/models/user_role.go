package models

import (
	"time"
)

type UserRole struct {
	UserID    int       `json:"user_id" gorm:"primaryKey"`
	RoleID    int       `json:"role_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Role Role `json:"role,omitempty" gorm:"foreignKey:RoleID;references:ID"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
