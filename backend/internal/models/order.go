package models

// OrderFlowerInfo представляет информацию о цветке в заказе (для JSON)
type OrderFlowerInfo struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}

// Модель заказа
type Order struct {
	Id         int               `json:"id" db:"id"`                // ID заказа
	UserId     int               `json:"userId" db:"user_id"`       // ID пользователя, который сделал заказ
	OrderDate  string            `json:"orderDate" db:"order_date"` // Дата
	TotalPrice string            `json:"totalPrice" db:"total_price"`
	Flowers    []OrderFlowerInfo `json:"flowers,omitempty"` // Добавляем поле для цветов
}

type OrderFlowers struct {
	OrderId  int `json:"orderId" db:"order_id"`
	FlowerId int `json:"flowerId" db:"flower_id" binding:"required"`
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
