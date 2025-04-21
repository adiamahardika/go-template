package models

type Todo struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
	DueDate     string `json:"due_date"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
