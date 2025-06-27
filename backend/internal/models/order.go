package models

// OrderFlowerInfo представляет информацию о цветке в заказе (для JSON)
type OrderFlowerInfo struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}

type Order struct {
	Id         int               `json:"id" db:"id"`
	UserId     int               `json:"user_id" db:"user_id"`
	OrderDate  string            `json:"order_date" db:"order_date"`
	TotalPrice string            `json:"total_price" db:"total_price"`
	Flowers    []OrderFlowerInfo `json:"flowers,omitempty"` // Добавляем поле для цветов
}

type OrderFlowers struct {
	OrderId  int `json:"order_id" db:"order_id"`
	FlowerId int `json:"flower_id" db:"flower_id" binding:"required"`
	Quantity int `json:"quantity" db:"quantity" binding:"required"`
}

type UpdateOrderInput struct {
	NewFlowerId int `json:"new_flower_id" binding:"required"`
	NewQuantity int `json:"new_quantity" binding:"required"`
}

type UpdateOrderFlowerIdInput struct {
	NewFlowerId int `json:"new_flower_id" binding:"required"`
}

type UpdateOrderQuantityInput struct {
	NewQuantity int `json:"new_quantity" binding:"required"`
}
