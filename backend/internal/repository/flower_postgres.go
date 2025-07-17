package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/christmas-fire/Bloomify/internal/apperror"
	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/jmoiron/sqlx"
)

type FlowerPostgres struct {
	db     *sqlx.DB
	logger *slog.Logger
}

// Параметры для фильтрации списка цветов
type FlowerFilter struct {
	Name        string   // Поиск по имени
	Description string   // Поиск по описанию
	MaxPrice    *float64 // Максимальная цена
	MaxStock    *int     // Максимальное количество
}

// Поля для обновления данных цветка
type UpdateFlowerInput struct {
	Name        *string  // Название
	Description *string  // Описание
	Price       *float64 // Цена
	Stock       *int     // Кол-во в наличии
}

func NewFlowerPostgres(db *sqlx.DB, logger *slog.Logger) *FlowerPostgres {
	return &FlowerPostgres{db: db, logger: logger}
}

func (r *FlowerPostgres) CreateFlower(ctx context.Context, name, description string, price float64, stock int) (int, error) {
	var id int
	query := "INSERT INTO flowers (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id"

	row := r.db.QueryRowxContext(ctx, query, name, description, price, stock)
	if err := row.Scan(&id); err != nil {
		// Обрабатываем ошибку duplicate key value violates unique constraint "\field\"
		if strings.Contains(err.Error(), "flowers_name_key") {
			return 0, fmt.Errorf("flower with name '%s' is already exists", name)
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (r *FlowerPostgres) Get(ctx context.Context, filter FlowerFilter) ([]models.Flower, error) {
	query := "SELECT id, name, description, price, stock FROM flowers WHERE 1=1"
	args := make([]any, 0)
	argId := 1

	if filter.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argId) // ILIKE для регистронезависимого поиска
		args = append(args, "%"+filter.Name+"%")           // % для поиска по подстроке
		argId++
	}

	if filter.Description != "" {
		query += fmt.Sprintf(" AND description ILIKE $%d", argId)
		args = append(args, "%"+filter.Description+"%")
		argId++
	}

	if filter.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", argId)
		args = append(args, *filter.MaxPrice)
		argId++
	}

	if filter.MaxStock != nil {
		query += fmt.Sprintf(" AND stock <= $%d", argId)
		args = append(args, *filter.MaxStock)
		argId++
	}

	query += " ORDER BY name ASC"

	flowers := make([]models.Flower, 0)
	if err := r.db.SelectContext(ctx, &flowers, query, args...); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return flowers, &apperror.TimeoutError{Err: err}
		}
		return flowers, &apperror.InternalServerError{Err: err}
	}
	return flowers, nil
}

func (r *FlowerPostgres) GetById(ctx context.Context, flowerId int) (models.Flower, error) {
	var flower models.Flower
	query := "SELECT id, name, description, price, stock FROM flowers WHERE id=$1"

	if err := r.db.GetContext(ctx, &flower, query, flowerId); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return flower, &apperror.TimeoutError{Err: err}
		}

		if errors.Is(err, sql.ErrNoRows) {
			return flower, &apperror.NotFoundError{
				Err:    err,
				Entity: "flower",
			}
		}
		return flower, &apperror.InternalServerError{Err: err}
	}
	return flower, nil
}

func (r *FlowerPostgres) Delete(ctx context.Context, flowerId int) error {
	query := "DELETE FROM flowers WHERE id=$1"

	if _, err := r.db.ExecContext(ctx, query, flowerId); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &apperror.TimeoutError{Err: err}
		}
		return &apperror.InternalServerError{Err: err}
	}
	return nil
}

func (r *FlowerPostgres) Update(ctx context.Context, id int, input UpdateFlowerInput) error {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argId := 1

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *input.Price)
		argId++
	}
	if input.Stock != nil {
		setValues = append(setValues, fmt.Sprintf("stock=$%d", argId))
		args = append(args, *input.Stock)
		argId++
	}

	if len(setValues) == 0 {
		return nil
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE flowers SET %s WHERE id=$%d", setQuery, argId)

	args = append(args, id)
	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return &apperror.TimeoutError{Err: err}
		}
		return &apperror.InternalServerError{Err: err}
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &apperror.InternalServerError{Err: err}
	}
	if rowsAffected == 0 {
		return &apperror.NotFoundError{
			Err:    nil,
			Entity: "flower",
		}
	}

	return nil
}
