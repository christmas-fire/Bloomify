package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/christmas-fire/Bloomify/config"
	_ "github.com/lib/pq"
)

const (
	SchemaUsers = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			jwt TEXT
		)`

	SchemaFlowers = `
		CREATE TABLE IF NOT EXISTS flowers (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			price NUMERIC NOT NULL,
			stock INT NOT NULL,
			image_url TEXT NOT NULL
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

var schemas = []string{SchemaUsers, SchemaFlowers, SchemaOrders, SchemaOrderFlowers}

func InitDB() *sql.DB {
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatal(err)
	}

	con := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		cfg.User, cfg.Password, cfg.Database, cfg.Host, cfg.Port, cfg.Sslmode,
	)

	db, err := sql.Open("postgres", con)
	if err != nil {
		log.Fatalf("error connect DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("DB unvaluable: %v", err)
	}

	log.Println("success connect to DB")

	return db
}

func InitTables(db *sql.DB) error {
	for _, schema := range schemas {
		_, err := db.Exec(schema)
		if err != nil {
			return fmt.Errorf("error exectuting schemas: %v", err)
		}
	}

	log.Println("success executing schemas")

	return nil
}
