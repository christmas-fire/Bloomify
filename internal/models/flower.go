package models

type Flower struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	ImageUrl string  `json:"image_url"`
}
