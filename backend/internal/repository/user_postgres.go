package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/christmas-fire/Bloomify/internal/apperror"
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

func (r *UserPostgres) GetAll(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0)
	query := "SELECT id, username, email, password FROM users"

	if err := r.db.SelectContext(ctx, &users, query); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return users, &apperror.TimeoutError{Err: err}
		}
		return users, &apperror.InternalServerError{Err: err}
	}
	return users, nil
}

func (r *UserPostgres) GetById(ctx context.Context, userId int) (models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password FROM users WHERE id=$1"

	if err := r.db.GetContext(ctx, &user, query, userId); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return user, &apperror.TimeoutError{Err: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			return user, &apperror.NotFoundError{
				Err:    err,
				Entity: "user",
			}
		}
		return user, &apperror.InternalServerError{Err: err}
	}
	return user, nil
}

func (r *UserPostgres) Delete(ctx context.Context, userId int) error {
	query := "DELETE FROM users WHERE id=$1"

	if _, err := r.db.ExecContext(ctx, query, userId); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &apperror.TimeoutError{Err: err}
		}
		return &apperror.InternalServerError{Err: err}
	}
	return nil
}

func (r *UserPostgres) UpdateUsername(ctx context.Context, userId int, oldUsername, newUsername string) error {
	var currentUsername string
	selectQuery := "SELECT username FROM users WHERE id=$1"

	if err := r.db.QueryRowContext(ctx, selectQuery, userId).Scan(&currentUsername); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &apperror.TimeoutError{Err: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			return &apperror.InvalidCredentialsError{
				Err: err,
			}
		}
		return &apperror.InternalServerError{Err: err}
	}

	if currentUsername != oldUsername {
		return &apperror.InvalidCredentialsError{
			Err: errors.New("incorrect username"),
		}
	}

	if newUsername == oldUsername {
		return &apperror.NoChangesError{
			Err: errors.New("username must be different"),
		}
	}

	updateQuery := "UPDATE users SET username=$1 WHERE id=$2"

	if _, err := r.db.ExecContext(ctx, updateQuery, newUsername, userId); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &apperror.TimeoutError{Err: err}
		}
		return &apperror.InternalServerError{Err: err}
	}
	return nil
}

func (r *UserPostgres) UpdatePassword(ctx context.Context, userId int, username, oldPassword, newPassword string) error {
	var currentHashedPassword string

	selectQuery := "SELECT password FROM users WHERE id=$1 AND username=$2"

	if err := r.db.QueryRowContext(ctx, selectQuery, userId, username).Scan(&currentHashedPassword); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &apperror.TimeoutError{Err: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			return &apperror.NotFoundError{
				Err:    err,
				Entity: "user",
			}
		}
		return &apperror.InternalServerError{Err: err}
	}

	if currentHashedPassword != oldPassword {
		return &apperror.InvalidCredentialsError{
			Err: errors.New("incorrect password"),
		}
	}

	if newPassword == oldPassword {
		return &apperror.NoChangesError{
			Err: errors.New("password must be different"),
		}
	}

	updateQuery := "UPDATE users SET password=$1 WHERE id=$2"

	if _, err := r.db.ExecContext(ctx, updateQuery, newPassword, userId); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &apperror.TimeoutError{Err: err}
		}
		return &apperror.InternalServerError{Err: err}
	}
	return nil
}
