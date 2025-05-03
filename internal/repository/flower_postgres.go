package repository

import (
	"database/sql"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/jmoiron/sqlx"
)

type FlowerPostgres struct {
	db *sqlx.DB
}

func NewFlowerPostgres(db *sqlx.DB) *FlowerPostgres {
	return &FlowerPostgres{db: db}
}

func (r *FlowerPostgres) CreateFlower(flower models.Flower) (int, error) {
	var id int
	query := "INSERT INTO flowers (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id"

	row := r.db.QueryRow(query, flower.Name, flower.Description, flower.Price, flower.Stock)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *FlowerPostgres) GetAll() ([]models.Flower, error) {
	var flowers []models.Flower
	query := "SELECT id, name, description, price, stock FROM flowers"

	err := r.db.Select(&flowers, query)

	if len(flowers) == 0 {
		return nil, nil
	}

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

func (r *FlowerPostgres) GetFlowersByName(name string) ([]models.Flower, error) {
	var flowers []models.Flower
	query := `SELECT * FROM flowers WHERE name ILIKE '%' || $1 || '%'`

	err := r.db.Select(&flowers, query, name)

	if err == sql.ErrNoRows {
		return []models.Flower{}, nil
	}
	return flowers, err
}

func (r *FlowerPostgres) GetFlowersByDescription(description string) ([]models.Flower, error) {
	var flowers []models.Flower
	query := "SELECT id, name, description, price, stock FROM flowers WHERE description ILIKE '%' || $1 || '%'"

	err := r.db.Select(&flowers, query, description)

	if len(flowers) == 0 {
		return nil, sql.ErrNoRows
	}

	return flowers, err
}

func (r *FlowerPostgres) GetFlowersByPrice(price float64) ([]models.Flower, error) {
	var flowers []models.Flower
	query := "SELECT id, name, description, price, stock FROM flowers WHERE price <= $1 ORDER BY price DESC"

	err := r.db.Select(&flowers, query, price)

	if len(flowers) == 0 {
		return nil, sql.ErrNoRows
	}

	return flowers, err
}

func (r *FlowerPostgres) GetFlowersByStock(stock int64) ([]models.Flower, error) {
	var flowers []models.Flower
	query := "SELECT id, name, description, price, stock FROM flowers WHERE stock <= $1 ORDER BY stock DESC"

	err := r.db.Select(&flowers, query, stock)

	if len(flowers) == 0 {
		return nil, sql.ErrNoRows
	}

	return flowers, err
}

func (r *FlowerPostgres) UpdateName(flowerId int, input models.UpdateNameInput) error {
	query := "UPDATE flowers SET name=$1 WHERE id=$2"

	_, err := r.db.Exec(query, input.NewName, flowerId)

	return err
}

func (r *FlowerPostgres) UpdateDescription(flowerId int, input models.UpdateDescriptionInput) error {
	query := "UPDATE flowers SET description=$1 WHERE id=$2"

	_, err := r.db.Exec(query, input.NewDescription, flowerId)

	return err
}

func (r *FlowerPostgres) UpdatePrice(flowerId int, input models.UpdatePriceInput) error {
	query := "UPDATE flowers SET price=$1 WHERE id=$2"

	_, err := r.db.Exec(query, input.NewPrice, flowerId)

	return err
}

func (r *FlowerPostgres) UpdateStock(flowerId int, input models.UpdateStockInput) error {
	query := "UPDATE flowers SET stock=$1 WHERE id=$2"

	_, err := r.db.Exec(query, input.NewStock, flowerId)

	return err
}
