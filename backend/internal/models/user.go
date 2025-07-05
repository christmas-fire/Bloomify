package models

// Модель пользователя
type User struct {
	Id       int    `json:"id" db:"id"`             // ID пользователя
	Username string `json:"username" db:"username"` // Имя пользователя
	Email    string `json:"email" db:"email"`       // Email пользователя
	Password string `json:"password" db:"password"` // Пароль пользователя
}
