package models

import (
	"time"
)

type PaymentMethod struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"not null"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:PaymentMethodID"`
}

func (PaymentMethod) TableName() string {
	return "payment_methods"
}
