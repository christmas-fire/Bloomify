package repository

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewAuthPostgres(db *sqlx.DB, logger *slog.Logger) *AuthPostgres {
	return &AuthPostgres{db: db, logger: logger}
}

func (r *AuthPostgres) CreateUser(username, email, password string) (int, error) {
	var id int
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"

	row := r.db.QueryRow(query, username, email, password)
	if err := row.Scan(&id); err != nil {
		// Обрабатываем ошибку duplicate key value violates unique constraint "\field\"
		if strings.Contains(err.Error(), "users_email_key") {
			return 0, fmt.Errorf("user with email '%s' is already exists", email)
		} else if strings.Contains(err.Error(), "users_username_key") {
			return 0, fmt.Errorf("user with username '%s' is already exists", username)
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE username=$1 AND password=$2"

	err := r.db.Get(&user, query, username, password)
	if err != nil {
		return user, err
	}

	return user, nil
}
