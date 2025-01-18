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
		return nil, sql.ErrNoRows
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
	query := "SELECT id, name, description, price, stock FROM flowers WHERE name=$1"

	err := r.db.Select(&flowers, query, name)

	if len(flowers) == 0 {
		return nil, sql.ErrNoRows
	}

	return flowers, err
}
