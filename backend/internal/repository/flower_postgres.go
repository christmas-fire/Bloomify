package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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

func (r *FlowerPostgres) CreateFlower(name, description string, price float64, stock int) (int, error) {
	var id int
	query := "INSERT INTO flowers (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id"

	row := r.db.QueryRow(query, name, description, price, stock)
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

func (r *FlowerPostgres) Get(filter FlowerFilter) ([]models.Flower, error) {
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

	var flowers []models.Flower
	err := r.db.Select(&flowers, query, args...)

	return flowers, err
}

func (r *FlowerPostgres) GetById(flowerId int) (models.Flower, error) {
	var flower models.Flower
	query := "SELECT id, name, description, price, stock FROM flowers WHERE id=$1"

	err := r.db.Get(&flower, query, flowerId)

	return flower, err
}

func (r *FlowerPostgres) Delete(flowerId int) error {
	query := "DELETE FROM flowers WHERE id=$1"

	_, err := r.db.Exec(query, flowerId)

	return err
}

func (r *FlowerPostgres) UpdateName(flowerId int, newName string) error {
	var currentName string
	selectQuery := "SELECT name FROM flowers WHERE id=$1"

	if err := r.db.QueryRow(selectQuery, flowerId).Scan(&currentName); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("flower not found")
		}
		return fmt.Errorf("failed to get current flower's name: %w", err)
	}

	if newName == currentName {
		return errors.New("you have no changes")
	}

	query := "UPDATE flowers SET name=$1 WHERE id=$2"
	_, err := r.db.Exec(query, newName, flowerId)

	return err
}

func (r *FlowerPostgres) UpdateDescription(flowerId int, newDescription string) error {
	var currentDescription string
	selectQuery := "SELECT description FROM flowers WHERE id=$1"

	if err := r.db.QueryRow(selectQuery, flowerId).Scan(&currentDescription); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("flower not found")
		}
		return fmt.Errorf("failed to get current flower's description: %w", err)
	}

	if newDescription == currentDescription {
		return errors.New("you have no changes")
	}

	query := "UPDATE flowers SET description=$1 WHERE id=$2"
	_, err := r.db.Exec(query, newDescription, flowerId)

	return err
}

func (r *FlowerPostgres) UpdatePrice(flowerId int, newPrice float64) error {
	var currentPrice float64
	selectQuery := "SELECT price FROM flowers WHERE id=$1"

	if err := r.db.QueryRow(selectQuery, flowerId).Scan(&currentPrice); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("flower not found")
		}
		return fmt.Errorf("failed to get current flower's price: %w", err)
	}

	if newPrice == currentPrice {
		return errors.New("you have no changes")
	}

	query := "UPDATE flowers SET price=$1 WHERE id=$2"
	_, err := r.db.Exec(query, newPrice, flowerId)

	return err
}

func (r *FlowerPostgres) UpdateStock(flowerId int, newStock int) error {
	var currentStock int
	selectQuery := "SELECT stock FROM flowers WHERE id=$1"

	if err := r.db.QueryRow(selectQuery, flowerId).Scan(&currentStock); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("flower not found")
		}
		return fmt.Errorf("failed to get current flower's stock: %w", err)
	}

	if newStock == currentStock {
		return errors.New("you have no changes")
	}

	query := "UPDATE flowers SET stock=$1 WHERE id=$2"
	_, err := r.db.Exec(query, newStock, flowerId)

	return err
}
