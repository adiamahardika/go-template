package models

import (
	"time"
)

type Product struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"not null"`
	Description *string    `json:"description,omitempty"`
	Price       float64    `json:"price" gorm:"type:decimal(12,2);not null"`
	Stock       int        `json:"stock" gorm:"default:0"`
	CategoryID  *int       `json:"category_id,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Category   *Category   `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
	CartItems  []CartItem  `json:"cart_items,omitempty" gorm:"foreignKey:ProductID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:ProductID"`
}

// DTO
type ProductResponse struct {
	ID           int               `json:"id"`
	Name         string            `json:"name"`
	Description  *string           `json:"description,omitempty"`
	Price        float64           `json:"price"`
	Stock        int               `json:"stock"`
	ImageURL     *string           `json:"image_url,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Category     *CategoryResponse `json:"category,omitempty"`
	Availability string            `json:"availability "`
	Related      []ProductRelated  `json:"related_product"`
}

// DTO
type ProductRelated struct {
	ID       int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string  `json:"name" gorm:"not null"`
	Price    float64 `json:"price" gorm:"type:decimal(12,2);not null"`
	ImageURL *string `json:"image_url,omitempty"`
}

// DTO
func (p *Product) ToProductResponse(r []ProductRelated) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Category: &CategoryResponse{
			ID:   p.Category.ID,
			Name: p.Category.Name,
		},
		Related: r,
	}
}

type CategoryResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (Product) TableName() string {
	return "products"
}
