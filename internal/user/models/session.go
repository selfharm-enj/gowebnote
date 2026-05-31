// Имплементация структуры Session.
package models

import "time"

// Структура для хранения Session.
type Session struct {
	ID           int64     `db:"id"`
	UserID       int64     `db:"user_id"`
	SessionToken string    `db:"session_token"`
	CSRFToken    string    `db:"csrf_token"`
	MaxAge       int       `db:"max_age"`
	ExpireAt     time.Time `db:"expired_at"`
}
