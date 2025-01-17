package models

// Модель цветка
type Flower struct {
	Id          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name" binding:"required"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price" binding:"required"`
	Stock       int     `json:"stock" db:"stock" binding:"required"`
}
