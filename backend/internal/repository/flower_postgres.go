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
	NameQuery        string   // Поиск по имени
	DescriptionQuery string   // Поиск по описанию
	MaxPrice         *float64 // Максимальная цена
	MaxStock         *int     // Максимальное количество
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
	args := []interface{}{}
	argId := 1

	if filter.NameQuery != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argId) // ILIKE для регистронезависимого поиска
		args = append(args, "%"+filter.NameQuery+"%")      // % для поиска по подстроке
		argId++
	}

	if filter.DescriptionQuery != "" {
		query += fmt.Sprintf(" AND description ILIKE $%d", argId)
		args = append(args, "%"+filter.DescriptionQuery+"%")
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

// TODO: объединить эти 4 обновления в один обработчик по типу как сделан Get()
func (r *FlowerPostgres) UpdateName(ctx context.Context, flowerId int, newName string) error {
	var currentName string
	selectQuery := "SELECT name FROM flowers WHERE id=$1"

	if err := r.db.QueryRowxContext(ctx, selectQuery, flowerId).Scan(&currentName); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("flower not found")
		}
		return fmt.Errorf("failed to get current flower's name: %w", err)
	}

	if newName == currentName {
		return errors.New("you have no changes")
	}

	query := "UPDATE flowers SET name=$1 WHERE id=$2"
	_, err := r.db.ExecContext(ctx, query, newName, flowerId)

	return err
}

func (r *FlowerPostgres) UpdateDescription(ctx context.Context, flowerId int, newDescription string) error {
	var currentDescription string
	selectQuery := "SELECT description FROM flowers WHERE id=$1"

	if err := r.db.QueryRowxContext(ctx, selectQuery, flowerId).Scan(&currentDescription); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("flower not found")
		}
		return fmt.Errorf("failed to get current flower's description: %w", err)
	}

	if newDescription == currentDescription {
		return errors.New("you have no changes")
	}

	query := "UPDATE flowers SET description=$1 WHERE id=$2"
	_, err := r.db.ExecContext(ctx, query, newDescription, flowerId)

	return err
}

func (r *FlowerPostgres) UpdatePrice(ctx context.Context, flowerId int, newPrice float64) error {
	var currentPrice float64
	selectQuery := "SELECT price FROM flowers WHERE id=$1"

	if err := r.db.QueryRowxContext(ctx, selectQuery, flowerId).Scan(&currentPrice); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("flower not found")
		}
		return fmt.Errorf("failed to get current flower's price: %w", err)
	}

	if newPrice == currentPrice {
		return errors.New("you have no changes")
	}

	query := "UPDATE flowers SET price=$1 WHERE id=$2"
	_, err := r.db.ExecContext(ctx, query, newPrice, flowerId)

	return err
}

func (r *FlowerPostgres) UpdateStock(ctx context.Context, flowerId int, newStock int) error {
	var currentStock int
	selectQuery := "SELECT stock FROM flowers WHERE id=$1"

	if err := r.db.QueryRowxContext(ctx, selectQuery, flowerId).Scan(&currentStock); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("flower not found")
		}
		return fmt.Errorf("failed to get current flower's stock: %w", err)
	}

	if newStock == currentStock {
		return errors.New("you have no changes")
	}

	query := "UPDATE flowers SET stock=$1 WHERE id=$2"
	_, err := r.db.ExecContext(ctx, query, newStock, flowerId)

	return err
}
