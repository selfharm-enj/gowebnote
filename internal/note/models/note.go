// Имплементация структуры Note.
package models

import "time"

// Структура для хранения Note.
type Note struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Name      string    `db:"name"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
}
