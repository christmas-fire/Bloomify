package models

// Модель цветка
type Flower struct {
	Id          int     `json:"id" db:"id"`                   // ID цветка
	Name        string  `json:"name" db:"name"`               // Название
	Description string  `json:"description" db:"description"` // Описание
	Price       float64 `json:"price" db:"price"`             // Цена
	Stock       int     `json:"stock" db:"stock"`             // Кол-во в наличии
}

type UpdateNameInput struct {
	NewName string `json:"newName" binding:"required"`
}

type UpdateDescriptionInput struct {
	NewDescription string `json:"newDescription" binding:"required"`
}

type UpdatePriceInput struct {
	NewPrice float64 `json:"newPrice" binding:"required"`
}

type UpdateStockInput struct {
	NewStock int `json:"newStock" binding:"required"`
}
