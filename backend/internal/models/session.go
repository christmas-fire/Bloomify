package models

import "time"

type RefreshSession struct {
	ID               int       `json:"-" db:"id"`
	UserID           int       `json:"-" db:"user_id"`
	RefreshTokenHash string    `json:"-" db:"refresh_token_hash"`
	ExpiresAt        time.Time `json:"-" db:"expires_at"`
	CreatedAt        time.Time `json:"-" db:"created_at"`
}
