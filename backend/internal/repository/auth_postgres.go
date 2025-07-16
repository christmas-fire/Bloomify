package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/christmas-fire/Bloomify/internal/apperror"
	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewAuthPostgres(db *sqlx.DB, logger *slog.Logger) *AuthPostgres {
	return &AuthPostgres{db: db, logger: logger}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, username, email, password string) (int, error) {
	var id int
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"

	row := r.db.QueryRowxContext(ctx, query, username, email, password)
	if err := row.Scan(&id); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return 0, &apperror.TimeoutError{Err: err}
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // 23505 	unique_violation
				if strings.Contains(pgErr.ConstraintName, "email") {
					return 0, &apperror.AlreadyExistsError{
						Err:    err,
						Entity: "user",
						Field:  "email",
						Value:  email,
					}
				}

				if strings.Contains(pgErr.ConstraintName, "username") {
					return 0, &apperror.AlreadyExistsError{
						Err:    err,
						Entity: "user",
						Field:  "username",
						Value:  username,
					}
				}
			}
		}
		return 0, &apperror.InternalServerError{Err: err}
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(ctx context.Context, username, password string) (models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE username=$1 AND password=$2"

	if err := r.db.GetContext(ctx, &user, query, username, password); err != nil {
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
