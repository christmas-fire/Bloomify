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

func (r *OrderPostgres) GetAllOrderFlowers() ([]models.OrderFlowers, error) {
	var orderFlowers []models.OrderFlowers
	query := "SELECT order_id, flower_id, quantity FROM order_flowers"

	err := r.db.Select(&orderFlowers, query)

	if len(orderFlowers) == 0 {
		return nil, sql.ErrNoRows
	}

	return orderFlowers, err
}

func (r *OrderPostgres) GetOrderFlowersByOrderId(orderFlowersId int) ([]models.OrderFlowers, error) {
	var orderFlowers []models.OrderFlowers
	query := "SELECT order_id, flower_id, quantity FROM order_flowers WHERE order_id=$1"

	err := r.db.Select(&orderFlowers, query, orderFlowersId)

	if len(orderFlowers) == 0 {
		return nil, sql.ErrNoRows
	}

	return orderFlowers, err
}

func (r *OrderPostgres) UpdateOrder(orderId int, input models.UpdateOrderInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Проверяем наличие цветка и получаем его цену
	var price float64
	query := "SELECT price FROM flowers WHERE id = $1"
	row := tx.QueryRow(query, input.NewFlowerId)
	err = row.Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("flower with id %d not found", input.NewFlowerId)
		}
		return err
	}

	// Устанавливаем новую общую стоимость заказа
	newTotalPrice := price * float64(input.NewQuantity)

	// Обновляем заказ
	query = `
        UPDATE orders 
        SET total_price = $1, order_date = NOW()
        WHERE id = $2
    `
	_, err = tx.Exec(query, newTotalPrice, orderId)
	if err != nil {
		return err
	}

	// Удаляем все старые записи о цветах в этом заказе
	query = "DELETE FROM order_flowers WHERE order_id = $1"
	_, err = tx.Exec(query, orderId)
	if err != nil {
		return err
	}

	// Добавляем новую запись о цветке в заказ
	query = "INSERT INTO order_flowers (order_id, flower_id, quantity) VALUES ($1, $2, $3)"
	_, err = tx.Exec(query, orderId, input.NewFlowerId, input.NewQuantity)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderPostgres) UpdateOrderFlowerId(orderId int, input models.UpdateOrderFlowerIdInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Проверяем наличие цветка и получаем его цену
	var newPrice float64
	query := "SELECT price FROM flowers WHERE id = $1"
	row := tx.QueryRow(query, input.NewFlowerId)
	err = row.Scan(&newPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("flower with id %d not found", input.NewFlowerId)
		}
		return err
	}

	// Получаем текущее количество цветков в заказе
	var quantity int
	query = "SELECT quantity FROM order_flowers WHERE order_id = $1"
	err = tx.QueryRow(query, orderId).Scan(&quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no flower found in order with id %d", orderId)
		}
		return err
	}

	// Устанавливаем новую общую стоимость заказа
	newTotalPrice := newPrice * float64(quantity)

	// Обновляем заказ
	query = `
        UPDATE orders 
        SET total_price = $1, order_date = NOW()
        WHERE id = $2
    `
	_, err = tx.Exec(query, newTotalPrice, orderId)
	if err != nil {
		return err
	}

	// Обновляем запись о цветке в заказе
	query = `
        UPDATE order_flowers 
        SET flower_id = $1
        WHERE order_id = $2
    `
	_, err = tx.Exec(query, input.NewFlowerId, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderPostgres) UpdateOrderQuantity(orderId int, input models.UpdateOrderQuantityInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Получаем текущий идентификатор цветка в заказе
	var flowerId int
	query := "SELECT flower_id FROM order_flowers WHERE order_id = $1"
	err = tx.QueryRow(query, orderId).Scan(&flowerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no flower found in order with id %d", orderId)
		}
		return err
	}

	// Проверяем наличие цветка и получаем его цену
	var price float64
	query = "SELECT price FROM flowers WHERE id = $1"
	row := tx.QueryRow(query, flowerId)
	err = row.Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("flower with id %d not found", flowerId)
		}
		return err
	}

	// Устанавливаем новую общую стоимость заказа
	newTotalPrice := price * float64(input.NewQuantity)

	// Обновляем заказ
	query = `
        UPDATE orders 
        SET total_price = $1, order_date = NOW()
        WHERE id = $2
    `
	_, err = tx.Exec(query, newTotalPrice, orderId)
	if err != nil {
		return err
	}

	// Обновляем запись о цветке в заказе
	query = `
        UPDATE order_flowers 
        SET quantity = $1
        WHERE order_id = $2
    `
	_, err = tx.Exec(query, input.NewQuantity, orderId)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderPostgres) Delete(orderId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Удаляем все записи о цветах в данном заказе
	query := "DELETE FROM order_flowers WHERE order_id = $1"
	result, err := tx.Exec(query, orderId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no flowers found in order with id %d", orderId)
	}

	// Удаляем сам заказ
	query = "DELETE FROM orders WHERE id = $1"
	result, err = tx.Exec(query, orderId)
	if err != nil {
		return err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("order with id %d not found", orderId)
	}

	return nil
}
