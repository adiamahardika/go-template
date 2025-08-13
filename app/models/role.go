package models

import (
	"time"
)

type Role struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"uniqueIndex;not null"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	UserRoles []UserRole `json:"user_roles,omitempty" gorm:"foreignKey:RoleID"`
}

func (Role) TableName() string {
	return "roles"
}
