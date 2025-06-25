package repository

import (
	"fmt"

	"github.com/christmas-fire/Bloomify/internal/models"
	"github.com/jmoiron/sqlx"
)

type SessionPostgres struct {
	db *sqlx.DB
}

func NewSessionPostgres(db *sqlx.DB) *SessionPostgres {
	return &SessionPostgres{db: db}
}

func (r *SessionPostgres) CreateSession(session models.RefreshSession) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, refresh_token_hash, expires_at) values ($1, $2, $3)", refreshSessionsTable)
	_, err := r.db.Exec(query, session.UserID, session.RefreshTokenHash, session.ExpiresAt)
	return err
}

func (r *SessionPostgres) GetSession(refreshTokenHash string) (models.RefreshSession, error) {
	var session models.RefreshSession
	query := fmt.Sprintf("SELECT id, user_id, refresh_token_hash, expires_at, created_at FROM %s WHERE refresh_token_hash=$1", refreshSessionsTable)
	err := r.db.Get(&session, query, refreshTokenHash)

	// Возвращаем sql.ErrNoRows как есть, если запись не найдена
	return session, err
}

func (r *SessionPostgres) DeleteSession(refreshTokenHash string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE refresh_token_hash=$1", refreshSessionsTable)
	_, err := r.db.Exec(query, refreshTokenHash)
	return err
}

// Добавляем константу для имени таблицы
const (
	refreshSessionsTable = "refresh_sessions"
)
