package repository

import (
	"database/sql"
	"fmt"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/jmoiron/sqlx"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) CreateOrder(userId int, orderFlower models.OrderFlowers) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var price float64
	query := "SELECT price FROM flowers WHERE id = $1"
	row := tx.QueryRow(query, orderFlower.FlowerId)
	err = row.Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("flower with id %d not found", orderFlower.FlowerId)
		}
		return 0, err
	}

	// Вычисляем общую стоимость заказа
	totalPrice := price * float64(orderFlower.Quantity)

	// Создаем новый заказ
	var orderId int
	query = "INSERT INTO orders (user_id, total_price, order_date) VALUES ($1, $2, NOW()) RETURNING id"
	err = tx.QueryRow(query, userId, totalPrice).Scan(&orderId)
	if err != nil {
		return 0, err
	}

	// Добавляем запись о цветке в заказ
	query = "INSERT INTO order_flowers (order_id, flower_id, quantity) VALUES ($1, $2, $3)"
	_, err = tx.Exec(query, orderId, orderFlower.FlowerId, orderFlower.Quantity)
	if err != nil {
		return 0, err
	}

	return orderId, nil
}

func (r *OrderPostgres) GetAll() ([]models.Order, error) {
	var orders []models.Order
	query := "SELECT id, user_id, order_date, total_price FROM orders"

	err := r.db.Select(&orders, query)

	if len(orders) == 0 {
		return nil, sql.ErrNoRows
	}

	return orders, err
}

func (r *OrderPostgres) GetById(orderId int) (models.Order, error) {
	var order models.Order
	query := "SELECT id, user_id, order_date, total_price FROM orders WHERE id=$1"

	err := r.db.Get(&order, query, orderId)

	return order, err
}

func (r *OrderPostgres) GetOrdersByUserId(userId int64) ([]models.Order, error) {
	var orders []models.Order
	query := "SELECT id, user_id, order_date, total_price FROM orders WHERE user_id=$1"

	err := r.db.Select(&orders, query, userId)

	if len(orders) == 0 {
		return nil, sql.ErrNoRows
	}

	return orders, err
}
