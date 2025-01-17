package models

// Модель пользователя
type User struct {
	Id       int    `json:"id" db:"id"`                                // ID пользователя
	Username string `json:"username" db:"username" binding:"required"` // Имя пользователя
	Email    string `json:"email" db:"email" binding:"required"`       // Email пользователя
	Password string `json:"password" db:"password" binding:"required"` // Пароль пользователя
}
