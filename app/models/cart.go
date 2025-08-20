package models

import (
	"time"
)

type Cart struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    *int      `json:"user_id,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CouponID  *int       `json:"coupon_id,omitempty"`

	// Relationships
	User      *User      `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
	CartItems []CartItem `json:"cart_items,omitempty" gorm:"foreignKey:CartID"`
	Coupon    *Coupon    `json:"coupon,omitempty" gorm:"foreignKey:CouponID;references:ID"`
}

func (Cart) TableName() string {
	return "carts"
}
