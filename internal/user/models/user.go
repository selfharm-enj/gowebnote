// Имплементация структуры User.
package models

import "time"

// Структура для хранения User.
type User struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
