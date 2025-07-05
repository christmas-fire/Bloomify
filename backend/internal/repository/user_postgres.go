package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewUserPostgres(db *sqlx.DB, logger *slog.Logger) *UserPostgres {
	return &UserPostgres{db: db, logger: logger}
}

func (r *UserPostgres) GetAll() ([]models.User, error) {
	var users []models.User
	query := "SELECT id, username, email, password FROM users"

	err := r.db.Select(&users, query)

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users, err
}

func (r *UserPostgres) GetById(userId int) (models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password FROM users WHERE id=$1"

	err := r.db.Get(&user, query, userId)

	return user, err
}

func (r *UserPostgres) Delete(userId int) error {
	query := "DELETE FROM users WHERE id=$1"

	_, err := r.db.Exec(query, userId)

	return err
}

func (r *UserPostgres) UpdateUsername(userId int, oldUsername, newUsername string) error {
	var currentUsername string
	selectQuery := "SELECT username FROM users WHERE id=$1"

	if err := r.db.QueryRow(selectQuery, userId).Scan(&currentUsername); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found or credentials do not match")
		}
		return fmt.Errorf("failed to get current username: %w", err)
	}

	if currentUsername != oldUsername {
		return errors.New("incorrect old username")
	}

	if newUsername == oldUsername {
		return errors.New("you have no changes")
	}

	updateQuery := "UPDATE users SET username=$1 WHERE id=$2"

	if _, err := r.db.Exec(updateQuery, newUsername, userId); err != nil {
		return fmt.Errorf("failed to update username: %w", err)

	}

	return nil
}

func (r *UserPostgres) UpdatePassword(userId int, username, oldPassword, newPassword string) error {
	var currentHashedPassword string

	selectQuery := "SELECT password FROM users WHERE id=$1 AND username=$2"
	if err := r.db.QueryRow(selectQuery, userId, username).Scan(&currentHashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found or credentials do not match")
		}
		return fmt.Errorf("failed to get current password: %w", err)
	}

	if currentHashedPassword != oldPassword {
		return errors.New("incorrect old password")
	}

	if newPassword == oldPassword {
		return errors.New("you have no changes")
	}

	updateQuery := "UPDATE users SET password=$1 WHERE id=$2"
	if _, err := r.db.Exec(updateQuery, newPassword, userId); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
