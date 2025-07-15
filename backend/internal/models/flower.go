package models

// Модель цветка
type Flower struct {
	Id          int     `json:"id" db:"id"`                   // ID цветка
	Name        string  `json:"name" db:"name"`               // Название
	Description string  `json:"description" db:"description"` // Описание
	Price       float64 `json:"price" db:"price"`             // Цена
	Stock       int     `json:"stock" db:"stock"`             // Кол-во в наличии
}
