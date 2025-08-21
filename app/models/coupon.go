package models

import (
	"time"
	"gorm.io/gorm"
)

type Coupon struct {
	ID              int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Code            string     `json:"code" gorm:"uniqueIndex;not null"`
	DiscountPercent *float64   `json:"discount_percent,omitempty" gorm:"type:decimal(5,2)"`
	MaxDiscount     *float64   `json:"max_discount,omitempty" gorm:"type:decimal(10,2)"`
	ExpiredAt       *time.Time `json:"expired_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Orders []Order `json:"orders,omitempty" gorm:"foreignKey:CouponID"`
}

func (Coupon) TableName() string {
	return "coupons"
}
