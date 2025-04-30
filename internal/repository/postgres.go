package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
			id SERIAL PRIMARY KEY unique,
			name TEXT NOT NULL unique,
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

	SchemaRefreshSessions = `
		CREATE TABLE IF NOT EXISTS refresh_sessions (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			refresh_token_hash VARCHAR(255) NOT NULL UNIQUE,
			expires_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	)`
)

// Схемы для инициализации базы данных
var schemas = []string{SchemaUsers, SchemaFlowers, SchemaOrders, SchemaOrderFlowers, SchemaRefreshSessions}

// Конфигурация PostgreSQL
type PostgreSQLConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Sslmode  string `yaml:"sslmode"`
}

// Инициализация PostgreSQL
func InitPostgreSQL(cfg PostgreSQLConfig) (*sqlx.DB, error) {
	con := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port, cfg.Sslmode,
	)

	db, err := sqlx.Open("postgres", con)
	if err != nil {
		return nil, fmt.Errorf("error open db: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connect db: %s", err.Error())
	}

	logrus.Println("success connect to db")

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
