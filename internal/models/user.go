package models

// Модель пользователя
type User struct {
	Id       int    `json:"id" db:"id"`                                // ID пользователя
	Username string `json:"username" db:"username" binding:"required"` // Имя пользователя
	Email    string `json:"email" db:"email" binding:"required"`       // Email пользователя
	Password string `json:"password" db:"password" binding:"required"` // Пароль пользователя
}

type UpdateUsernameInput struct {
	OldUsername string `json:"old_username" binding:"required"`
	NewUsername string `json:"new_username" binding:"required"`
}

type UpdatePasswordInput struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
