package models

import (
	"time"
)

type Category struct {
	ID          int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"not null"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	// Relationships
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

type CategoryAdmin struct {
	ID             int        `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string     `json:"name" gorm:"not null"`
	Description    *string    `json:"description,omitempty"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	IncludeRelated bool       `json:"include_related,omitempty" gorm:"index"`
	IncludeDeleted bool       `json:"include_deleted,omitempty" gorm:"index"`
	// Relationships
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

type CategoryResponses struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Products []Product `json:"products, omitempty"`
}

type CategoryInfo struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description`
}

// DTO
func (c *Category) ToCategoryResponse(r []Product) CategoryResponses {
	return CategoryResponses{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		Products:    r,
	}
}

func (Category) TableName() string {
	return "categories"
}
