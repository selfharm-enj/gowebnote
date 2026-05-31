// Вспомогательные типы для NoteService.
package usecase

import (
	"errors"
	"time"
)

// Формат запроса на создание Note.
type CreateNoteReq struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

// Валидация NotesReq.
func (r *CreateNoteReq) Validate() error {
	if len(r.Name) == 0 || r.Name == "" {
		return errors.New("name are empty")
	}
	if len(r.Text) == 0 || r.Text == "" {
		return errors.New("text are empty")
	}
	return nil
}

type NoteIDRes struct {
	ID int64 `json:"id"`
}

// Формат ответа Note.
type NoteRes struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type NoteReq struct {
	ID int64 `json:"id"`
}

func (r *NoteReq) Validate() error {
	if r.ID <= 0 {
		return errors.New("invalid ID value")
	}
	return nil
}

type UpdateNoteReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Text string `json:"text"`
}

func (r *UpdateNoteReq) Validate() error {
	if r.ID <= 0 {
		return errors.New("invalid ID value")
	}
	if len(r.Name) == 0 || r.Name == "" {
		return errors.New("name are empty")
	}
	if len(r.Text) == 0 || r.Text == "" {
		return errors.New("text are empty")
	}
	return nil
}
