package models

import (
	"time"
)

type User struct {
	ID        int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string     `json:"name" gorm:"not null"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Password  string     `json:"password" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Carts     []Cart     `json:"carts,omitempty" gorm:"foreignKey:UserID"`
	Orders    []Order    `json:"orders,omitempty" gorm:"foreignKey:UserID"`
	UserRoles []UserRole `json:"user_roles,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
