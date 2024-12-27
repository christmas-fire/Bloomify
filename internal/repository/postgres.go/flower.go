package postgres

import (
	"database/sql"
	"fmt"

	"github.com/christmas-fire/Bloomify/internal/models"
)

// Репозиторий для работы с цветами
type FlowerRepository struct {
	db *sql.DB
}

// Создание нового репозитория для работы с цветами
func NewFlowerRepository(db *sql.DB) *FlowerRepository {
	return &FlowerRepository{db: db}
}

// Валидация цветка
func ValidateFlower(f models.Flower) error {
	if f.Name == "" {
		return fmt.Errorf("name can't be empty")
	}
	if f.Price <= 0 {
		return fmt.Errorf("price can't be less or equal 0")
	}
	if f.Stock < 0 {
		return fmt.Errorf("stock can't negative number")
	}
	return nil
}

// Добавление цветка
func (r *FlowerRepository) AddFlower(f models.Flower) error {
	checkQuery := `
		SELECT EXISTS (
			SELECT 1
			FROM flowers
			WHERE name = $1
	)`

	insertQuery := `
		INSERT INTO flowers (name, price, stock)
		VALUES ($1, $2, $3)
	`

	if err := ValidateFlower(f); err != nil {
		return fmt.Errorf("invalid flower data: %w", err)
	}

	var exists bool
	if err := r.db.QueryRow(checkQuery, f.Name).Scan(&exists); err != nil {
		return fmt.Errorf("error check if flower exists: %w", err)
	}

	if exists {
		return fmt.Errorf("flower already exists")
	}

	if _, err := r.db.Exec(insertQuery, f.Name, f.Price, f.Stock); err != nil {
		return fmt.Errorf("error inserting new flower into database: %w", err)
	}

	return nil
}

// Получение всех цветов
func (r *FlowerRepository) GetAllFlowers() ([]models.Flower, error) {
	query := `
		SELECT id, name, price, stock FROM flowers
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flowers []models.Flower
	for rows.Next() {
		var f models.Flower
		if err := rows.Scan(&f.Id, &f.Name, &f.Price, &f.Stock); err != nil {
			return nil, err
		}
		flowers = append(flowers, f)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return flowers, nil
}

// Получение цветка по ID
func (r *FlowerRepository) GetFlowerByID(id int) (*models.Flower, error) {
	query := `
		SELECT id, name, price, stock
		FROM flowers
		WHERE id = $1
	`

	var f models.Flower
	err := r.db.QueryRow(query, id).Scan(&f.Id, &f.Name, &f.Price, &f.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("flower with id %d not found", id)
		}
		return nil, fmt.Errorf("error querying flower by id: %w", err)
	}

	return &f, nil
}

// Получение цветов по названию
func (r *FlowerRepository) GetFlowersByName(f models.Flower) ([]models.Flower, error) {
    query := `
        SELECT id, name, price, stock
        FROM flowers
        WHERE name ILIKE $1
    `
    
    rows, err := r.db.Query(query, "%"+f.Name+"%")
    if err != nil {
        return nil, fmt.Errorf("error querying flowers by name: %w", err)
    }
    defer rows.Close()

    var flowers []models.Flower
    for rows.Next() {
        var flower models.Flower
        if err := rows.Scan(&flower.Id, &flower.Name, &flower.Price, &flower.Stock); err != nil {
            return nil, fmt.Errorf("error scanning flower row: %w", err)
        }
        flowers = append(flowers, flower)
    }

    return flowers, nil
}

// Получение цветов по цене
func (r *FlowerRepository) GetFlowersByPrice(f models.Flower) ([]models.Flower, error) {
    query := `
        SELECT id, name, price, stock
        FROM flowers
        WHERE price <= $1
        ORDER BY price ASC
    `
    
    rows, err := r.db.Query(query, f.Price)
    if err != nil {
        return nil, fmt.Errorf("error querying flowers by price: %w", err)
    }
    defer rows.Close()

    var flowers []models.Flower
    for rows.Next() {
        var flower models.Flower
        if err := rows.Scan(&flower.Id, &flower.Name, &flower.Price, &flower.Stock); err != nil {
            return nil, fmt.Errorf("error scanning flower row: %w", err)
        }
        flowers = append(flowers, flower)
    }

    return flowers, nil
}

// Получение цветов по количеству в наличии
func (r *FlowerRepository) GetFlowersByStock(f models.Flower) ([]models.Flower, error) {
    query := `
        SELECT id, name, price, stock
        FROM flowers
        WHERE stock >= $1
        ORDER BY stock DESC
    `
    
    rows, err := r.db.Query(query, f.Stock)
    if err != nil {
        return nil, fmt.Errorf("error querying flowers by stock: %w", err)
    }
    defer rows.Close()

    var flowers []models.Flower
    for rows.Next() {
        var flower models.Flower
        if err := rows.Scan(&flower.Id, &flower.Name, &flower.Price, &flower.Stock); err != nil {
            return nil, fmt.Errorf("error scanning flower row: %w", err)
        }
        flowers = append(flowers, flower)
    }

    return flowers, nil
}

// Удаление цветка по ID
func (r *FlowerRepository) DeleteFlowerByID(id int) error {
	query := `
		DELETE FROM flowers
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting flower with id '%d': %v", id, err)
	}

	return nil
}
