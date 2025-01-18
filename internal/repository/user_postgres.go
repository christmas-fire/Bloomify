package repository

import (
	"database/sql"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
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

func (r *UserPostgres) UpdateUsername(userId int, input models.UpdateUsernameInput) error {
	query := "UPDATE users SET username=$1 WHERE username=$2 AND id=$3"

	_, err := r.db.Exec(query, input.NewUsername, input.OldUsername, userId)

	return err
}

func (r *UserPostgres) UpdatePassword(userId int, input models.UpdatePasswordInput) error {
	query := "UPDATE users SET password=$1 WHERE username=$2 AND password=$3"

	_, err := r.db.Exec(query, input.NewPassword, input.Username, input.OldPassword)

	return err
}
