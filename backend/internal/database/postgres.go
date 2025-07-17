package database

import (
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

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
