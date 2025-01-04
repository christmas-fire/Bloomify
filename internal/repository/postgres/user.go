package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/christmas-fire/Bloomify/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// Репозиторий для работы с пользователями
type UserRepository struct {
	db *sql.DB
}

// Создание нового репозитория для работы с пользователями
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Валидация пользователя
func ValidateUser(u models.User) error {
	if len(u.Username) < 3 {
		return fmt.Errorf("username must have at least 3 characters")
	}

	if len(u.Password) < 8 {
		return fmt.Errorf("password must have at least 8 characters")
	}

	if !strings.Contains(u.Email, "@") {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// Регистрация пользователя
func (r *UserRepository) Register(u models.User) error {
	if err := ValidateUser(u); err != nil {
		return fmt.Errorf("invalid user data: %w", err)
	}

	checkQuery := `
		SELECT EXISTS (
			SELECT 1
			FROM USERS
			WHERE username = $1 OR email = $2
		)
	`

	var exists bool
	if err := r.db.QueryRow(checkQuery, u.Username, u.Email).Scan(&exists); err != nil {
		return fmt.Errorf("error check if user exists: %w", err)
	}

	if exists {
		return fmt.Errorf("user with the same username/email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	insertQuery := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
	`

	if _, err := r.db.Exec(insertQuery, u.Username, u.Email, hashedPassword); err != nil {
		return fmt.Errorf("error inserting new user into database: %w", err)
	}

	return nil
}

// Вход пользователя
func (r *UserRepository) Login(u models.User) error {
	query := `
		SELECT password FROM users WHERE username = $1
	`

	var hashedPassword string
	if err := r.db.QueryRow(query, u.Username).Scan(&hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("error retrieving user from the database: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(u.Password)); err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}

// Получение всех пользователей
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `
		SELECT id, username, email, password FROM users
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.Username, &u.Email, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Удаление пользователя по ID
func (r *UserRepository) DeleteUserByID(id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user with id '%d': %v", id, err)
	}

	return nil
}
