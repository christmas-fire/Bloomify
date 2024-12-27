package models

// Модель пользователя
type User struct {
	Id       int    `json:"id"` // ID пользователя
	Username string `json:"username"` // Имя пользователя
	Email    string `json:"email"` // Email пользователя
	Password string `json:"password"` // Пароль пользователя
}
