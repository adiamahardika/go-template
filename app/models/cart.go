package models

import (
	"time"
)

type Cart struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    *int      `json:"user_id,omitempty"`
	CouponID  *int      `json:"coupon_id,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	User   *User      `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	Coupon *Coupon    `json:"coupon,omitempty" gorm:"foreignKey:CouponID;references:ID"`
	Items  []CartItem `json:"items,omitempty" gorm:"foreignKey:CartID"`
}

func (Cart) TableName() string {
	return "carts"
}
