package database

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Константы для схем базы данных
const (
	SchemaUsers = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)`

	SchemaFlowers = `
		CREATE TABLE IF NOT EXISTS flowers (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			price NUMERIC NOT NULL,
			stock INT NOT NULL
		)`

	SchemaOrders = `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			order_date TIMESTAMP DEFAULT NOW(),
			total_price NUMERIC NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
		)`

	SchemaOrderFlowers = `
		CREATE TABLE IF NOT EXISTS order_flowers (
			order_id INT NOT NULL,
			flower_id INT NOT NULL,
			quantity INT NOT NULL,
			PRIMARY KEY (order_id, flower_id),
			FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
			FOREIGN KEY (flower_id) REFERENCES flowers (id) ON DELETE CASCADE
		)`
)

// Схемы для инициализации базы данных
var schemas = []string{SchemaUsers, SchemaFlowers, SchemaOrders, SchemaOrderFlowers}

// Инициализация Postgres
func InitPostgres() (*sqlx.DB, error) {
	// Получаем строку подключения из .env
	dsn := os.Getenv("POSTGRES_DSN")

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		logrus.Fatal(err)
	}

	return db, nil
}

// Инициализация таблиц
func InitTables(db *sqlx.DB) error {
	for _, schema := range schemas {
		_, err := db.Exec(schema)
		if err != nil {
			return fmt.Errorf("error exectuting schema: %v", err)
		}
	}

	logrus.Println("success executing schemas")

	return nil
}
