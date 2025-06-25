package repository

import (
	"database/sql"
	"fmt"
	"strconv"

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
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 1. Найти активный (последний) orderId для пользователя
	var orderId int
	var currentTotalPriceStr string
	queryOrder := "SELECT id, total_price FROM orders WHERE user_id=$1 ORDER BY order_date DESC LIMIT 1"
	err = tx.QueryRow(queryOrder, userId).Scan(&orderId, &currentTotalPriceStr)

	orderExists := true
	if err != nil {
		if err == sql.ErrNoRows {
			orderExists = false // Заказа нет, будем создавать новый
		} else {
			return 0, fmt.Errorf("failed to find order for user %d: %w", userId, err)
		}
	}

	// 2. Получить цену и остаток добавляемого цветка
	var price float64
	var stock int                                                            // Добавляем переменную для остатка
	queryPrice := "SELECT price, stock FROM flowers WHERE id = $1"           // Обновляем запрос
	err = tx.QueryRow(queryPrice, orderFlower.FlowerId).Scan(&price, &stock) // Сканируем цену и остаток
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("flower with id %d not found", orderFlower.FlowerId)
		}
		return 0, fmt.Errorf("failed to get price and stock for flower %d: %w", orderFlower.FlowerId, err)
	}

	itemTotalPrice := price * float64(orderFlower.Quantity)
	newOrderTotalPrice := itemTotalPrice // Для нового заказа

	if orderExists {
		// Заказ уже есть, проверяем, есть ли там такой цветок
		var existingQuantity int
		queryCheckFlower := "SELECT quantity FROM order_flowers WHERE order_id=$1 AND flower_id=$2"
		err = tx.QueryRow(queryCheckFlower, orderId, orderFlower.FlowerId).Scan(&existingQuantity)

		if err != nil && err != sql.ErrNoRows {
			return 0, fmt.Errorf("failed to check flower %d in order %d: %w", orderFlower.FlowerId, orderId, err)
		}

		// Парсим текущую общую цену заказа
		currentTotalPrice, parseErr := strconv.ParseFloat(currentTotalPriceStr, 64)
		if parseErr != nil {
			return 0, fmt.Errorf("failed to parse current total price for order %d: %w", orderId, parseErr)
		}

		if err == sql.ErrNoRows {
			// Проверка на stock перед добавлением
			if orderFlower.Quantity > stock {
				return 0, fmt.Errorf("insufficient stock for flower %d. Available: %d, Requested: %d", orderFlower.FlowerId, stock, orderFlower.Quantity)
			}
			// Цветка в заказе нет - добавляем новую строку в order_flowers
			queryInsertOF := "INSERT INTO order_flowers (order_id, flower_id, quantity) VALUES ($1, $2, $3)"
			_, err = tx.Exec(queryInsertOF, orderId, orderFlower.FlowerId, orderFlower.Quantity)
			if err != nil {
				return 0, fmt.Errorf("failed to insert flower %d into order %d: %w", orderFlower.FlowerId, orderId, err)
			}
			newOrderTotalPrice = currentTotalPrice + itemTotalPrice // Обновляем общую цену
		} else {
			// Проверка на stock перед обновлением
			if existingQuantity+orderFlower.Quantity > stock {
				return 0, fmt.Errorf("insufficient stock for flower %d. Available: %d, In cart: %d, Requested to add: %d", orderFlower.FlowerId, stock, existingQuantity, orderFlower.Quantity)
			}
			// Цветок в заказе есть - обновляем количество
			newQuantity := existingQuantity + orderFlower.Quantity
			queryUpdateOF := "UPDATE order_flowers SET quantity = $1 WHERE order_id = $2 AND flower_id = $3"
			_, err = tx.Exec(queryUpdateOF, newQuantity, orderId, orderFlower.FlowerId)
			if err != nil {
				return 0, fmt.Errorf("failed to update quantity for flower %d in order %d: %w", orderFlower.FlowerId, orderId, err)
			}
			newOrderTotalPrice = currentTotalPrice + itemTotalPrice // Обновляем общую цену
		}

		// Обновляем total_price в существующем заказе
		newTotalPriceStr := fmt.Sprintf("%.2f", newOrderTotalPrice)
		queryUpdateOrder := "UPDATE orders SET total_price = $1, order_date = NOW() WHERE id = $2"
		_, err = tx.Exec(queryUpdateOrder, newTotalPriceStr, orderId)
		if err != nil {
			return 0, fmt.Errorf("failed to update total price for order %d: %w", orderId, err)
		}

	} else {
		// Проверка на stock перед созданием заказа с этим товаром
		if orderFlower.Quantity > stock {
			return 0, fmt.Errorf("insufficient stock for flower %d to create order. Available: %d, Requested: %d", orderFlower.FlowerId, stock, orderFlower.Quantity)
		}
		// Заказа нет - создаем новый заказ и добавляем первую строку в order_flowers
		newTotalPriceStr := fmt.Sprintf("%.2f", newOrderTotalPrice)
		queryNewOrder := "INSERT INTO orders (user_id, total_price, order_date) VALUES ($1, $2, NOW()) RETURNING id"
		err = tx.QueryRow(queryNewOrder, userId, newTotalPriceStr).Scan(&orderId)
		if err != nil {
			return 0, fmt.Errorf("failed to create new order for user %d: %w", userId, err)
		}

		queryInsertOF := "INSERT INTO order_flowers (order_id, flower_id, quantity) VALUES ($1, $2, $3)"
		_, err = tx.Exec(queryInsertOF, orderId, orderFlower.FlowerId, orderFlower.Quantity)
		if err != nil {
			return 0, fmt.Errorf("failed to insert flower %d into new order %d: %w", orderFlower.FlowerId, orderId, err)
		}
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
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.Order{}, nil
		}
		return nil, err
	}

	for i := range orders {
		var flowers []models.OrderFlowerInfo
		flowerQuery := "SELECT flower_id as id, quantity FROM order_flowers WHERE order_id=$1"
		err = r.db.Select(&flowers, flowerQuery, orders[i].Id)
		if err != nil && err != sql.ErrNoRows {
			return nil, fmt.Errorf("error fetching flowers for order %d: %w", orders[i].Id, err)
		}
		orders[i].Flowers = flowers
	}

	return orders, nil
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

	newTotalPrice := price * float64(input.NewQuantity)

	query = `
        UPDATE orders 
        SET total_price = $1, order_date = NOW()
        WHERE id = $2
    `
	_, err = tx.Exec(query, newTotalPrice, orderId)
	if err != nil {
		return err
	}

	query = "DELETE FROM order_flowers WHERE order_id = $1"
	_, err = tx.Exec(query, orderId)
	if err != nil {
		return err
	}

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

	var quantity int
	query = "SELECT quantity FROM order_flowers WHERE order_id = $1"
	err = tx.QueryRow(query, orderId).Scan(&quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no flower found in order with id %d", orderId)
		}
		return err
	}

	newTotalPrice := newPrice * float64(quantity)

	query = `
        UPDATE orders 
        SET total_price = $1, order_date = NOW()
        WHERE id = $2
    `
	_, err = tx.Exec(query, newTotalPrice, orderId)
	if err != nil {
		return err
	}

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

	var flowerId int
	query := "SELECT flower_id FROM order_flowers WHERE order_id = $1"
	err = tx.QueryRow(query, orderId).Scan(&flowerId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no flower found in order with id %d", orderId)
		}
		return err
	}

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

	newTotalPrice := price * float64(input.NewQuantity)

	query = `
        UPDATE orders 
        SET total_price = $1, order_date = NOW()
        WHERE id = $2
    `
	_, err = tx.Exec(query, newTotalPrice, orderId)
	if err != nil {
		return err
	}

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

// RemoveFlowerFromOrderByUser удаляет указанный цветок из активного заказа пользователя.
// Предполагается, что у пользователя только один активный заказ (корзина).
func (r *OrderPostgres) RemoveFlowerFromOrderByUser(userId int, flowerId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	// 1. Найти активный orderId для userId (предполагаем, что он один, берем последний)
	var orderId int
	queryOrder := "SELECT id FROM orders WHERE user_id=$1 ORDER BY order_date DESC LIMIT 1"
	err = tx.QueryRow(queryOrder, userId).Scan(&orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no active order found for user %d", userId)
		}
		return fmt.Errorf("failed to find order for user %d: %w", userId, err)
	}

	// 2. Получить цену удаляемого цветка и его количество в заказе
	var price float64
	var quantity int
	queryFlowerInfo := `
		SELECT f.price, of.quantity 
		FROM order_flowers of 
		JOIN flowers f ON of.flower_id = f.id 
		WHERE of.order_id = $1 AND of.flower_id = $2`
	err = tx.QueryRow(queryFlowerInfo, orderId, flowerId).Scan(&price, &quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("flower %d not found in order %d", flowerId, orderId)
		}
		return fmt.Errorf("failed to get flower info for order %d, flower %d: %w", orderId, flowerId, err)
	}

	// 3. Удалить запись из order_flowers
	queryDelete := "DELETE FROM order_flowers WHERE order_id = $1 AND flower_id = $2"
	result, err := tx.Exec(queryDelete, orderId, flowerId)
	if err != nil {
		return fmt.Errorf("failed to delete flower %d from order %d: %w", flowerId, orderId, err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// Это не должно произойти, если предыдущий запрос нашел цветок, но проверим
		return fmt.Errorf("flower %d was not found in order %d during delete", flowerId, orderId)
	}

	// 4. Обновить total_price в заказе
	var currentTotalPriceStr string
	queryGetPrice := "SELECT total_price FROM orders WHERE id = $1"
	err = tx.QueryRow(queryGetPrice, orderId).Scan(&currentTotalPriceStr)
	if err != nil {
		return fmt.Errorf("failed to get current total price for order %d: %w", orderId, err)
	}
	currentTotalPrice, err := strconv.ParseFloat(currentTotalPriceStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse current total price '%s' for order %d: %w", currentTotalPriceStr, orderId, err)
	}

	priceReduction := price * float64(quantity)
	newTotalPrice := currentTotalPrice - priceReduction
	newTotalPriceStr := fmt.Sprintf("%.2f", newTotalPrice) // Форматируем до 2 знаков

	queryUpdateOrder := "UPDATE orders SET total_price = $1 WHERE id = $2"
	_, err = tx.Exec(queryUpdateOrder, newTotalPriceStr, orderId)
	if err != nil {
		return fmt.Errorf("failed to update total price for order %d: %w", orderId, err)
	}

	return nil
}

// IncrementFlowerQuantity увеличивает количество цветка в активном заказе пользователя на 1
func (r *OrderPostgres) IncrementFlowerQuantity(userId int, flowerId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 1. Найти активный orderId
	var orderId int
	queryOrder := "SELECT id FROM orders WHERE user_id=$1 ORDER BY order_date DESC LIMIT 1"
	err = tx.QueryRow(queryOrder, userId).Scan(&orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no active order found for user %d", userId)
		}
		return fmt.Errorf("failed to find order for user %d: %w", userId, err)
	}

	// 2. Получить цену цветка и остаток на складе
	var price float64
	var stock int
	queryFlower := "SELECT price, stock FROM flowers WHERE id = $1"
	err = tx.QueryRow(queryFlower, flowerId).Scan(&price, &stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("flower with id %d not found", flowerId)
		}
		return fmt.Errorf("failed to get price and stock for flower %d: %w", flowerId, err)
	}

	// 3. Проверить текущее количество в корзине и сравнить с остатком
	var currentQuantity int
	queryQuantity := "SELECT quantity FROM order_flowers WHERE order_id = $1 AND flower_id = $2"
	err = tx.QueryRow(queryQuantity, orderId, flowerId).Scan(&currentQuantity)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если строки нет, это значит, что цветок не в корзине, инкрементировать нечего.
			// Либо это ошибка, т.к. инкремент вызывается для существующего в корзине товара.
			return fmt.Errorf("flower %d not found in order %d for increment", flowerId, orderId)
		}
		return fmt.Errorf("failed to get current quantity for order %d, flower %d: %w", orderId, flowerId, err)
	}

	// Проверка на stock перед инкрементом
	if currentQuantity+1 > stock {
		return fmt.Errorf("insufficient stock for flower %d. Available: %d, In cart: %d", flowerId, stock, currentQuantity)
	}

	// 4. Увеличить quantity в order_flowers
	queryUpdateOF := "UPDATE order_flowers SET quantity = quantity + 1 WHERE order_id = $1 AND flower_id = $2"
	result, err := tx.Exec(queryUpdateOF, orderId, flowerId)
	if err != nil {
		return fmt.Errorf("failed to increment quantity for order %d, flower %d: %w", orderId, flowerId, err)
	}
	rowsAffected, _ := result.RowsAffected() // Проверяем, была ли строка обновлена
	if rowsAffected == 0 {
		// Эта проверка дублирует проверку выше (на sql.ErrNoRows), но оставим для надежности
		return fmt.Errorf("flower %d not found in order %d during increment update", flowerId, orderId)
	}

	// 5. Обновить total_price в orders
	var currentTotalPriceStr string
	queryGetPriceInc := "SELECT total_price FROM orders WHERE id = $1"
	err = tx.QueryRow(queryGetPriceInc, orderId).Scan(&currentTotalPriceStr)
	if err != nil {
		return fmt.Errorf("failed to get current total price for order %d: %w", orderId, err)
	}
	currentTotalPriceInc, err := strconv.ParseFloat(currentTotalPriceStr, 64)
	if err != nil {
		return fmt.Errorf("failed to parse current total price '%s' for order %d: %w", currentTotalPriceStr, orderId, err)
	}

	newTotalPriceInc := currentTotalPriceInc + price
	newTotalPriceStrInc := fmt.Sprintf("%.2f", newTotalPriceInc)

	queryUpdateOrderInc := "UPDATE orders SET total_price = $1 WHERE id = $2"
	_, err = tx.Exec(queryUpdateOrderInc, newTotalPriceStrInc, orderId)
	if err != nil {
		return fmt.Errorf("failed to update total price for order %d: %w", orderId, err)
	}

	return nil
}

// DecrementFlowerQuantity уменьшает количество цветка. Если = 0, удаляет.
// TODO: Вернуть кастомную ошибку, если цветок удален.
func (r *OrderPostgres) DecrementFlowerQuantity(userId int, flowerId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 1. Найти активный orderId
	var orderId int
	queryOrder := "SELECT id FROM orders WHERE user_id=$1 ORDER BY order_date DESC LIMIT 1"
	err = tx.QueryRow(queryOrder, userId).Scan(&orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no active order found for user %d", userId)
		}
		return fmt.Errorf("failed to find order for user %d: %w", userId, err)
	}

	// 2. Получить текущее количество и цену цветка
	var price float64
	var currentQuantity int
	queryFlowerInfo := `
		SELECT f.price, of.quantity 
		FROM order_flowers of 
		JOIN flowers f ON of.flower_id = f.id 
		WHERE of.order_id = $1 AND of.flower_id = $2`
	err = tx.QueryRow(queryFlowerInfo, orderId, flowerId).Scan(&price, &currentQuantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("flower %d not found in order %d", flowerId, orderId)
		}
		return fmt.Errorf("failed to get flower info for order %d, flower %d: %w", orderId, flowerId, err)
	}

	if currentQuantity <= 1 {
		// Удаляем цветок, если количество 1 или меньше
		queryDelete := "DELETE FROM order_flowers WHERE order_id = $1 AND flower_id = $2"
		_, err = tx.Exec(queryDelete, orderId, flowerId)
		if err != nil {
			return fmt.Errorf("failed to delete flower %d from order %d (quantity <= 1): %w", flowerId, orderId, err)
		}
		// TODO: Возвращать кастомную ошибку? service.ErrFlowerRemoved
	} else {
		// Уменьшаем количество на 1
		queryUpdateOF := "UPDATE order_flowers SET quantity = quantity - 1 WHERE order_id = $1 AND flower_id = $2"
		_, err = tx.Exec(queryUpdateOF, orderId, flowerId)
		if err != nil {
			return fmt.Errorf("failed to decrement quantity for order %d, flower %d: %w", orderId, flowerId, err)
		}
	}

	// 4. Обновить total_price в orders (уменьшаем на цену одного цветка)
	var currentTotalPriceStrDec string
	queryGetPriceDec := "SELECT total_price FROM orders WHERE id = $1"
	err = tx.QueryRow(queryGetPriceDec, orderId).Scan(&currentTotalPriceStrDec)
	if err != nil {
		return fmt.Errorf("failed to get current total price for order %d: %w", orderId, err)
	}
	currentTotalPriceDec, err := strconv.ParseFloat(currentTotalPriceStrDec, 64)
	if err != nil {
		return fmt.Errorf("failed to parse current total price '%s' for order %d: %w", currentTotalPriceStrDec, orderId, err)
	}

	newTotalPriceDec := currentTotalPriceDec - price
	newTotalPriceStrDec := fmt.Sprintf("%.2f", newTotalPriceDec)

	queryUpdateOrderDec := "UPDATE orders SET total_price = $1 WHERE id = $2"
	_, err = tx.Exec(queryUpdateOrderDec, newTotalPriceStrDec, orderId)
	if err != nil {
		return fmt.Errorf("failed to update total price for order %d: %w", orderId, err)
	}

	return nil
}

// DeleteActiveOrderByUserId находит и удаляет последний (активный) заказ пользователя.
func (r *OrderPostgres) DeleteActiveOrderByUserId(userId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 1. Найти ID последнего заказа пользователя
	var orderId int
	queryFindOrder := "SELECT id FROM orders WHERE user_id = $1 ORDER BY order_date DESC LIMIT 1"
	err = tx.QueryRow(queryFindOrder, userId).Scan(&orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			// Активного заказа нет, нечего удалять
			return nil // Не считаем это ошибкой
		}
		return fmt.Errorf("failed to find active order for user %d: %w", userId, err)
	}

	// 2. Удалить связанные записи из order_flowers
	queryDeleteFlowers := "DELETE FROM order_flowers WHERE order_id = $1"
	_, err = tx.Exec(queryDeleteFlowers, orderId)
	if err != nil {
		// Не фатально, если записей не было, но логируем
		// logrus.Warnf("Could not delete order_flowers for order %d (maybe none existed?): %v", orderId, err)
		// Но если ошибка другая - возвращаем
		// Чтобы точно знать, можно было бы проверить RowsAffected, но пока пропустим
		// Проверка нужна т.к. Exec не возвращает sql.ErrNoRows
		// Простая проверка:
		var exists int // Проверим, был ли хоть один цветок
		_ = tx.QueryRow("SELECT 1 FROM order_flowers WHERE order_id = $1 LIMIT 1", orderId).Scan(&exists)
		if err != nil && exists == 1 { // Если ошибка И цветок был, то проблема
			return fmt.Errorf("failed to delete order_flowers for order %d: %w", orderId, err)
		}
		// Если ошибки нет или цветка не было, продолжаем
		err = nil // Сбрасываем ошибку, если она была некритичной
	}

	// 3. Удалить сам заказ из orders
	queryDeleteOrder := "DELETE FROM orders WHERE id = $1"
	result, err := tx.Exec(queryDeleteOrder, orderId)
	if err != nil {
		return fmt.Errorf("failed to delete order %d: %w", orderId, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// Это не должно произойти, если мы нашли orderId на шаге 1
		return fmt.Errorf("active order %d for user %d disappeared during transaction", orderId, userId)
	}

	return nil // Успех
}
