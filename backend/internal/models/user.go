package models

// Модель пользователя
type User struct {
	Id       int    `json:"id" db:"id"`                                // ID пользователя
	Username string `json:"username" db:"username" binding:"required"` // Имя пользователя
	Email    string `json:"email" db:"email" binding:"required"`       // Email пользователя
	Password string `json:"password" db:"password" binding:"required"` // Пароль пользователя
}

type UpdateUsernameInput struct {
	OldUsername string `json:"oldUsername" binding:"required"`
	NewUsername string `json:"newUsername" binding:"required"`
}

type UpdatePasswordInput struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}
