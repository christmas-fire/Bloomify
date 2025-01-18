package models

// Модель цветка
type Flower struct {
	Id          int     `json:"id" db:"id"`
	Name        string  `json:"name" db:"name" binding:"required"`
	Description string  `json:"description" db:"description"`
	Price       float64 `json:"price" db:"price" binding:"required"`
	Stock       int     `json:"stock" db:"stock" binding:"required"`
}

type UpdateNameInput struct {
	NewName string `json:"new_name" binding:"required"`
}

type UpdateDescriptionInput struct {
	NewDescription string `json:"new_description" binding:"required"`
}

type UpdatePriceInput struct {
	NewPrice float64 `json:"new_price" binding:"required"`
}

type UpdateStockInput struct {
	NewStock int `json:"new_stock" binding:"required"`
}
