// Интерфейс Note.
package repository

import (
	"context"

	"webnote/internal/note/models"
)

// Интерфейс для работы с Note.
type NoteRepository interface {
	CreateNote(ctx context.Context, note models.Note) (int64, error)
	GetNote(ctx context.Context, noteID int64, userID int64) (models.Note, error)
	GetMyNotes(ctx context.Context, userID int64) ([]models.Note, error)
	UpdateNote(ctx context.Context, noteID int64, userID int64, newNote models.Note) (models.Note, error)
	DeleteNote(ctx context.Context, noteID int64, userID int64) (int64, error)
}
