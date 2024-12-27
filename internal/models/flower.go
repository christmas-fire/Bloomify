package models

// Модель цветка
type Flower struct {
	Id       int     `json:"id"` // ID цветка
	Name     string  `json:"name"` // Название цветка
	Price    float64 `json:"price"` // Цена цветка
	Stock    int     `json:"stock"` // Количество в наличии
	ImageUrl string  `json:"image_url"` // URL изображения цветка
}
